package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/gorilla/mux"
)

type SeasonHandler struct {
	repo    db.Repository
	planner *domain.PredictivePlanner
}

func NewSeasonHandler(repo db.Repository) *SeasonHandler {
	return &SeasonHandler{
		repo:    repo,
		planner: domain.NewPredictivePlanner(repo),
	}
}

func (h *SeasonHandler) HandlePredictiveLoadout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	historicalShowID, err := strconv.ParseInt(vars["historical_show_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid historical show ID", http.StatusBadRequest)
		return
	}

	predictedRings, err := h.planner.PredictShowLoadout(r.Context(), historicalShowID)
	if err != nil {
		http.Error(w, "Failed to generate prediction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictedRings)
}

func (h *SeasonHandler) HandleApplyPredictedLoadout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newShowID, err := strconv.ParseInt(vars["show_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid new show ID", http.StatusBadRequest)
		return
	}

	var rings []domain.ShowRing
	if err := json.NewDecoder(r.Body).Decode(&rings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.planner.ApplyPredictedLoadout(r.Context(), newShowID, rings); err != nil {
		http.Error(w, "Failed to apply loadout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *SeasonHandler) HandleGetShowsForSeason(w http.ResponseWriter, r *http.Request) {
	// Dummy implementation for brevity, typically would query repo.GetShowsForSeason
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]domain.Show{})
}
