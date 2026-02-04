package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/desmond/rental-management-system/cmd/server/docs" // Setup for Swagger docs
	"github.com/desmond/rental-management-system/internal/api"
	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/fleet"
	_ "github.com/lib/pq"
)

// @title Rental Management System API
// @version 1.0
// @description API for managing fleet assets, rentals, and maintenance.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// Remote Management Setup
	registry := fleet.NewRemoteRegistry()
	mockMgr := fleet.NewMockRemoteManager()
	registry.Register("mock-provider", mockMgr)

	handler := api.NewHandler(repo, registry)
	router := api.NewRouter(handler)

	// Swagger UI
	// We need to type assert router to *http.ServeMux to add the swagger handler,
	// or perform this inside api.NewRouter.
	// To keep it simple, we'll assume NewRouter returns a Handler but we can verify.
	// Actually, let's just modify the router inside NewRouter, but since that returns http.Handler,
	// we will wrap the router or cast it if we know the implementation.
	// BETTER: Add it to the router inside api.NewRouter, but we need the httpSwagger import there.
	// ALTERNATIVE: Use a default mux here for swagger if not interfering.

	// Just add it to the standard mux in router.go, much cleaner.

	log.Printf("Starting server on :8080 (DB: %s)", dbURL)
	if err := http.ListenAndServe(":8080", router); err != nil {

		log.Fatal(err)
	}
}
