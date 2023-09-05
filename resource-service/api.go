package main

import "net/http"

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) handleCreateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resource Resource
		err := s.readJSON(w, r, &resource)
		if err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		err = s.svc.createResource(&resource)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, http.StatusOK, resource)
	}
}

func (s *ApiServer) handleGetResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.readIDParam(r)
		if err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		resource, err := s.svc.getResource(id)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, http.StatusOK, resource)
	}
}

func (s *ApiServer) handleGetResources() http.HandlerFunc {
	type input struct {
		Name       string
		Jobtitles  []string
		Types      []string
		Workgroups []string
		Locations  []string
		Managers   []string
		Filters
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var input input
		v := New()
		qs := r.URL.Query()
		input.Name = s.readString(qs, "name", "")
		input.Jobtitles = s.readCSV(qs, "titles", []string{})
		input.Types = s.readCSV(qs, "types", []string{})
		input.Workgroups = s.readCSV(qs, "workgroups", []string{})
		input.Locations = s.readCSV(qs, "locations", []string{})
		input.Managers = s.readCSV(qs, "managers", []string{})
		input.Filters.Page = s.readInt(qs, "page", 1, v)
		input.Filters.PageSize = s.readInt(qs, "pagesize", 20, v)
		input.Filters.Sort = s.readString(qs, "sort", "id")
		input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}
		if ValidateFilters(v, input.Filters); !v.Valid() {
			s.failedValidation(w, r, http.StatusBadRequest, v.Errors)
			return
		}
		resources, metadata, err := s.svc.getResources(input.Name, input.Jobtitles, input.Types, input.Workgroups, input.Locations, input.Managers, input.Filters)
		if err != nil {
			// TODO: Amend to check for ErrNoRows in case there are not rows in the resources table
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		err = s.writeJSON(w, http.StatusOK, envelope{"resources": resources, "metadata": metadata}, nil)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}

func (s *ApiServer) handleUpdateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resource Resource
		err := s.readJSON(w, r, &resource)
		if err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		err = s.svc.updateResource(&resource)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		s.writeJSON(w, http.StatusOK, resource)
	}
}
