package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/desmond/rental-management-system/internal/domain"
)

// CreateScheduledDelivery handles the creation of a new scheduled delivery.
func (h *Handler) CreateScheduledDelivery(w http.ResponseWriter, r *http.Request) {
	var delivery domain.ScheduledDelivery
	if err := json.NewDecoder(r.Body).Decode(&delivery); err != nil {
		log.Printf("failed to decode CreateScheduledDelivery request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateScheduledDelivery(r.Context(), &delivery); err != nil {
		log.Printf("failed to create scheduled delivery: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(delivery)
}

// GetScheduledDelivery handles retrieving a scheduled delivery by ID.
func (h *Handler) GetScheduledDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/deliveries/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid delivery id", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.GetScheduledDeliveryByID(r.Context(), id)
	if err != nil {
		log.Printf("failed to get scheduled delivery: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if delivery == nil {
		http.Error(w, "delivery not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(delivery)
}

// ListScheduledDeliveries handles retrieving a list of scheduled deliveries, optionally filtered by event ID.
func (h *Handler) ListScheduledDeliveries(w http.ResponseWriter, r *http.Request) {
	var eventID *int64
	if eIDStr := r.URL.Query().Get("event_id"); eIDStr != "" {
		if eID, err := strconv.ParseInt(eIDStr, 10, 64); err == nil {
			eventID = &eID
		}
	}

	deliveries, err := h.repo.ListScheduledDeliveries(r.Context(), eventID)
	if err != nil {
		log.Printf("failed to list scheduled deliveries: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deliveries)
}

// CreateShipment handles the creation of a new shipment.
func (h *Handler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var shipment domain.Shipment
	if err := json.NewDecoder(r.Body).Decode(&shipment); err != nil {
		log.Printf("failed to decode CreateShipment request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateShipment(r.Context(), &shipment); err != nil {
		log.Printf("failed to create shipment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shipment)
}

// GetShipment handles retrieving a shipment by ID.
func (h *Handler) GetShipment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/shipments/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid shipment id", http.StatusBadRequest)
		return
	}

	shipment, err := h.repo.GetShipmentByID(r.Context(), id)
	if err != nil {
		log.Printf("failed to get shipment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shipment == nil {
		http.Error(w, "shipment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(shipment)
}

// ListShipments handles retrieving a list of shipments, optionally filtered by delivery ID.
func (h *Handler) ListShipments(w http.ResponseWriter, r *http.Request) {
	var deliveryID *int64
	if dIDStr := r.URL.Query().Get("delivery_id"); dIDStr != "" {
		if dID, err := strconv.ParseInt(dIDStr, 10, 64); err == nil {
			deliveryID = &dID
		}
	}

	shipments, err := h.repo.ListShipments(r.Context(), deliveryID)
	if err != nil {
		log.Printf("failed to list shipments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(shipments)
}

// UpdateShipment handles updating an existing shipment.
func (h *Handler) UpdateShipment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/shipments/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid shipment id", http.StatusBadRequest)
		return
	}

	var shipment domain.Shipment
	if err := json.NewDecoder(r.Body).Decode(&shipment); err != nil {
		log.Printf("failed to decode UpdateShipment request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	shipment.ID = id

	if err := h.repo.UpdateShipment(r.Context(), &shipment); err != nil {
		log.Printf("failed to update shipment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(shipment)
}

// AllocateAssets handles bulk allocation of assets to a shipment.
func (h *Handler) AllocateAssets(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/logistics/shipments/")
	idStr = strings.TrimSuffix(idStr, "/allocate")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid shipment id", http.StatusBadRequest)
		return
	}

	var req struct {
		AssetIDs []int64 `json:"asset_ids"`
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

	if err := h.repo.AllocateAssetsToShipment(r.Context(), id, req.AssetIDs, *agentIDVal); err != nil {
		log.Printf("failed to allocate assets to shipment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
