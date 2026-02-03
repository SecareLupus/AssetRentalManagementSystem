package main

import (
	"log"
	"net/http"

	"github.com/desmond/rental-management-system/internal/api"
)

func main() {
	handler := api.NewHandler()
	router := api.NewRouter(handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
