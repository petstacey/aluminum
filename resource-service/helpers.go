package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/petstacey/iter"
	"github.com/petstacey/validator"
)

type envelope map[string]any

func (a *ApiServer) readIDParam(r *http.Request) (int64, error) {
	param := iter.Param(r.Context(), "id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (a *ApiServer) readString(qs url.Values, key string, defaultValue string) string {
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}
	return str
}

func (a *ApiServer) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func (a *ApiServer) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		v.AddError(key, "must be 'true' or 'false'")
		return defaultValue
	}
	return val
}

func (a *ApiServer) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var UnmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &UnmarshalTypeError):
			if UnmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", UnmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", UnmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

func (a *ApiServer) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (a *ApiServer) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	return a.writeJSON(w, statusCode, envelope{"error": err.Error()})
}

func (a *ApiServer) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := a.writeJSON(w, status, envelope{"error": message}, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *ApiServer) failedValidation(w http.ResponseWriter, r *http.Request, status int, errors map[string]string) {
	a.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
