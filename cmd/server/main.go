package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/desmond/rental-management-system/cmd/server/docs" // Setup for Swagger docs
	"github.com/desmond/rental-management-system/internal/api"
	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/desmond/rental-management-system/internal/fleet"
	"github.com/desmond/rental-management-system/internal/mqtt"
	"github.com/desmond/rental-management-system/internal/worker"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// Automated Migrations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.RunMigrations(ctx, conn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := db.NewSqlRepository(conn)

	// Seed Admin User
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser != "" && adminPass != "" {
		ctx := context.Background()
		existing, err := repo.GetUserByUsername(ctx, adminUser)
		if err != nil {
			log.Printf("Error checking for admin user: %v", err)
		} else if existing == nil {
			log.Printf("Creating default admin user: %s", adminUser)
			hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed to hash admin password: %v", err)
			} else {
				user := &domain.User{
					Username:     adminUser,
					PasswordHash: string(hash),
					Email:        "admin@example.com", // Default email
					Role:         domain.UserRoleAdmin,
					IsEnabled:    true,
				}
				if err := repo.CreateUser(ctx, user); err != nil {
					log.Printf("Failed to create admin user: %v", err)
				} else {
					log.Printf("Admin user created successfully")
				}
			}
		}
	}

	// MQTT Setup
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "tcp://localhost:1883"
	}
	mqttClient := mqtt.NewClient(mqtt.Config{
		BrokerURL: mqttBroker,
		ClientID:  "rms-server",
	})
	if err := mqttClient.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to MQTT broker: %v", err)
	} else {
		defer mqttClient.Disconnect()
	}

	// Remote Management Setup
	registry := fleet.NewRemoteRegistry()
	mockMgr := fleet.NewMockRemoteManager()
	registry.Register("mock-provider", mockMgr)

	// Workers
	outboxWorker := worker.NewOutboxWorker(repo, mqttClient)
	go outboxWorker.Start(context.Background(), 5*time.Second)

	healthWorker := worker.NewHealthWorker(repo, mqttClient, registry)
	go healthWorker.Start(context.Background(), 1*time.Minute)

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
