package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/desmond/rental-management-system/internal/domain"
)

// GetRentalFulfillment retrieves the fulfillment details for a reservation.
// @Summary Get Fulfillment Status
// @Description Calculates the delta between demands and actual asset movements.
// @Tags Logistics
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} domain.RentalFulfillmentStatus
// @Router /logistics/reservations/{id}/fulfillment [get]
func (h *Handler) GetRentalFulfillment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/reservations/")
	idStr = strings.TrimSuffix(idStr, "/fulfillment")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	status, err := h.repo.GetRentalFulfillmentStatus(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// BatchDispatchAssets performs a batch checkout of assets for a reservation.
// @Summary Batch Dispatch Assets
// @Description Records multiple CheckOutActions and transitions assets to 'deployed'.
// @Tags Logistics
// @Accept json
// @Param id path int true "Reservation ID"
// @Param request body object{asset_ids=[]int64,to_location_id=int64} true "Dispatch Data"
// @Success 204 {string} string "No Content"
// @Router /logistics/reservations/{id}/dispatch [post]
func (h *Handler) BatchDispatchAssets(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/reservations/")
	idStr = strings.TrimSuffix(idStr, "/dispatch")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		AssetIDs       []int64 `json:"asset_ids"`
		FromLocationID *int64  `json:"from_location_id"`
		ToLocationID   *int64  `json:"to_location_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	agentIDVal := h.getUserIDFromContext(r)
	if agentIDVal == nil {
		http.Error(w, "agent id missing from context", http.StatusUnauthorized)
		return
	}

	if err := h.repo.BatchCheckOut(r.Context(), id, req.AssetIDs, *agentIDVal, req.FromLocationID, req.ToLocationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append Outbox Event for Batch Dispatch
	payload, _ := json.Marshal(map[string]interface{}{
		"reservation_id": id,
		"asset_count":    len(req.AssetIDs),
		"agent_id":       agentIDVal,
	})
	h.repo.AppendEvent(r.Context(), nil, &domain.OutboxEvent{
		Type:    domain.EventAssetTransitioned,
		Payload: payload,
	})

	w.WriteHeader(http.StatusNoContent)
}

// BatchReturnAssets performs a batch return of assets.
// @Summary Batch Return Assets
// @Description Records multiple ReturnActions and transitions assets to 'available'.
// @Tags Logistics
// @Accept json
// @Param id path int true "Reservation ID"
// @Param request body object{asset_ids=[]int64} true "Return Data"
// @Success 204 {string} string "No Content"
// @Router /logistics/reservations/{id}/return [post]
func (h *Handler) BatchReturnAssets(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/reservations/")
	idStr = strings.TrimSuffix(idStr, "/return")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		AssetIDs     []int64 `json:"asset_ids"`
		ToLocationID *int64  `json:"to_location_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	agentIDVal := h.getUserIDFromContext(r)
	if agentIDVal == nil {
		http.Error(w, "agent id missing from context", http.StatusUnauthorized)
		return
	}

	if err := h.repo.BatchReturn(r.Context(), id, req.AssetIDs, *agentIDVal, req.ToLocationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
