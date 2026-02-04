package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/desmond/rental-management-system/internal/api"
	"github.com/desmond/rental-management-system/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repo := db.NewSqlRepository(conn)
	handler := api.NewHandler(repo)
	router := api.NewRouter(handler)

	log.Printf("Starting server on :8080 (DB: %s)", dbURL)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
