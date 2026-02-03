package api

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	// db db.Repository
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "GetCatalog stub"})
}

func (h *Handler) CreateRentAction(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "CreateRentAction stub"})
}

func (h *Handler) GetRentAction(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "GetRentAction stub"})
}
