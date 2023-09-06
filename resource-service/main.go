package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/petstacey/aluminum/resource-service/data"
	"github.com/petstacey/iter"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	DSN = "host=localhost port=5432 user=postgres password=password dbname=aluminum sslmode=disable timezone=UTC connect_timeout=5"
)

func main() {
	db := connectToDB()
	repo := data.NewPostgresRepository(db)
	svc := NewService(repo)
	svc = NewLoggingService(svc)
	api := NewApiServer(svc)
	rtr := iter.New()
	rtr.HandlerFunc("/v1/resources", api.handleCreateResource(), http.MethodPost)
	rtr.HandlerFunc("/v1/resources", api.handleGetResources(), http.MethodGet)
	rtr.HandlerFunc("/v1/resources", api.handleUpdateResource(), http.MethodPut)
	rtr.HandlerFunc("/v1/resources/:id", api.handleGetResource(), http.MethodGet)
	http.ListenAndServe(":6543", rtr)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	counts := 0
	for {
		conn, err := openDB(DSN)
		if err != nil {
			fmt.Printf("Postgres not yet ready: %v\n", err)
			counts++
		} else {
			fmt.Println("Connected to PostgreSQL!")
			return conn
		}
		if counts > 10 {
			fmt.Println("Retry limit reached and not connected to database")
			return nil
		}
		fmt.Println("Backing off and waiting for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
