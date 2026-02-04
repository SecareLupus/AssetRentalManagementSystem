package api

import (
	"encoding/json"
	"net/http"

	"github.com/desmond/rental-management-system/internal/db"
)

type Handler struct {
	repo db.Repository
}

func NewHandler(repo db.Repository) *Handler {
	return &Handler{repo: repo}
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
