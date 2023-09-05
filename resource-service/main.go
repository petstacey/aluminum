package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	DSN = "host=localhost port=5432 user=postgres password=password dbname=aluminum sslmode=disable timezone=UTC connect_timeout=5"
)

func main() {
	db := connectToDB()
	repo := NewPostgresRepository(db)
	svc := NewService(repo)
	svc = NewLoggingService(svc)
	// resource := Resource{
	// 	ID:        234567,
	// 	Name:      "Jill Doe",
	// 	Email:     "jill@doe.com",
	// 	Type:      "Full time",
	// 	JobTitle:  "Senior Manager",
	// 	Workgroup: "Managers and Non-billable - ANZ",
	// 	Location:  "ACT",
	// 	Manager:   "Jane Doe",
	// }
	// err := svc.createResource(&resource)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }
	// resources, _, err := svc.getResources("", []string{}, []string{}, []string{}, []string{}, []string{}, Filters{Page: 1, PageSize: 20, Sort: "id", SortSafelist: []string{"id"}})
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// }
	// for _, resource := range resources {
	// 	fmt.Printf("%d: %s\n", resource.ID, resource.Name)
	// }
	api := NewApiServer(svc)
	rtr := httprouter.New()
	rtr.HandlerFunc(http.MethodPost, "/v1/resources", api.handleCreateResource())
	rtr.HandlerFunc(http.MethodGet, "/v1/resources", api.handleGetResources())
	rtr.HandlerFunc(http.MethodPut, "/v1/resources", api.handleUpdateResource())
	rtr.HandlerFunc(http.MethodGet, "/v1/resources/:id", api.handleGetResource())
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
			fmt.Println("Count not connect to database")
			return nil
		}
		fmt.Println("Backing off and waiting for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
