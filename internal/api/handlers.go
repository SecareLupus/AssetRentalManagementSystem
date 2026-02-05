package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/desmond/rental-management-system/internal/fleet"
)

type Handler struct {
	repo           db.Repository
	remoteRegistry *fleet.RemoteRegistry
}

func NewHandler(repo db.Repository, remoteRegistry *fleet.RemoteRegistry) *Handler {
	return &Handler{
		repo:           repo,
		remoteRegistry: remoteRegistry,
	}
}

func (h *Handler) getUserIDFromContext(r *http.Request) *int64 {
	claims, ok := r.Context().Value(UserContextKey).(map[string]interface{}) // JWT usually unmarshals to map[string]interface{} or jwt.MapClaims
	if !ok {
		// Try jwt.MapClaims
		if c, ok := r.Context().Value(UserContextKey).(map[string]interface{}); ok {
			claims = c
		} else {
			// Depending on how jwt runs, it might be distinct type.
			// In auth.go: ctx := context.WithValue(r.Context(), UserContextKey, claims) where claims is jwt.MapClaims
			// claims is map[string]interface{}
			return nil
		}
	}
	// jwt.MapClaims is map[string]interface{}
	claims, ok = r.Context().Value(UserContextKey).(map[string]interface{})
	if !ok {
		return nil
	}

	if idFloat, ok := claims["user_id"].(float64); ok {
		id := int64(idFloat)
		return &id
	}
	return nil
}

// ItemType Handlers

func (h *Handler) validateItemType(it *domain.ItemType) error {
	if it.Code == "" {
		return fmt.Errorf("code is required")
	}
	if it.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch it.Kind {
	case domain.ItemKindSerialized, domain.ItemKindFungible, domain.ItemKindKit:
		// valid
	default:
		return fmt.Errorf("invalid kind: %s", it.Kind)
	}
	return nil
}

// CreateItemType creates a new item type in the catalog.
// @Summary Create Item Type
// @Description Creates a new Item Type definition.
// @Tags Catalog
// @Accept json
// @Produce json
// @Param item_type body domain.ItemType true "Item Type Definition"
// @Success 201 {object} domain.ItemType
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /catalog/item-types [post]
func (h *Handler) CreateItemType(w http.ResponseWriter, r *http.Request) {
	var it domain.ItemType
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validateItemType(&it); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateItemType(r.Context(), &it); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}

func (h *Handler) GetItemType(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/catalog/item-types/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	it, err := h.repo.GetItemTypeByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if it == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(it)
}

func (h *Handler) UpdateItemType(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/catalog/item-types/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var it domain.ItemType
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	it.ID = id

	if err := h.validateItemType(&it); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateItemType(r.Context(), &it); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(it)
}

func (h *Handler) DeleteItemType(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/catalog/item-types/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteItemType(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCatalog returns all item types.
// @Summary List Item Types
// @Description Returns the full catalog of Item Types.
// @Tags Catalog
// @Produce json
// @Success 200 {array} domain.ItemType
// @Failure 500 {string} string "Internal Server Error"
// @Router /catalog/item-types [get]
func (h *Handler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	results, err := h.repo.ListItemTypes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Asset Handlers

func (h *Handler) validateAsset(a *domain.Asset) error {
	if a.ItemTypeID == 0 {
		return fmt.Errorf("item_type_id is required")
	}
	switch a.Status {
	case domain.AssetStatusAvailable, domain.AssetStatusReserved, domain.AssetStatusMaintenance, domain.AssetStatusRetired:
		// valid
	case "":
		a.Status = domain.AssetStatusAvailable
	default:
		return fmt.Errorf("invalid status: %s", a.Status)
	}
	return nil
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var a domain.Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validateAsset(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.CreatedByUserID = h.getUserIDFromContext(r)

	if err := h.repo.CreateAsset(r.Context(), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

// GetAsset retrieves a specific asset.
// @Summary Get Asset
// @Description Retrieves an asset by its ID.
// @Tags Assets
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} domain.Asset
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventory/assets/{id} [get]
func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	a, err := h.repo.GetAssetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if a == nil {
		http.NotFound(w, r)
		return
	}

	// Fetch ItemType to check features
	it, err := h.repo.GetItemTypeByID(r.Context(), a.ItemTypeID)
	if err == nil && it != nil {
		if !it.SupportedFeatures.RemoteManagement {
			a.RemoteManagementID = nil
		}
		if !it.SupportedFeatures.Provisioning {
			a.ProvisioningStatus = ""
			a.FirmwareVersion = nil
			a.Hostname = nil
		}
		if !it.SupportedFeatures.BuildSpecTracking {
			a.BuildSpecVersion = nil
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

// UpdateAsset updates an existing asset.
// @Summary Update Asset
// @Description Updates an existing asset's details.
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param asset body domain.Asset true "Asset Data"
// @Success 200 {object} domain.Asset
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventory/assets/{id} [put]
func (h *Handler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var a domain.Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	a.ID = id

	if err := h.validateAsset(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch ItemType to check features before saving
	it, err := h.repo.GetItemTypeByID(r.Context(), a.ItemTypeID)
	if err == nil && it != nil {
		if !it.SupportedFeatures.RemoteManagement {
			a.RemoteManagementID = nil
		}
		if !it.SupportedFeatures.Provisioning {
			a.ProvisioningStatus = ""
			a.FirmwareVersion = nil
			a.Hostname = nil
		}
		if !it.SupportedFeatures.BuildSpecTracking {
			a.BuildSpecVersion = nil
		}
	}

	if err := h.repo.UpdateAsset(r.Context(), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

// UpdateAssetStatus updates the status of an asset.
// @Summary Update Asset Status
// @Description Updates the status of an asset (e.g., available, maintenance).
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param status body object{status=domain.AssetStatus} true "New Status"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventory/assets/{id}/status [patch]
func (h *Handler) UpdateAssetStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Status domain.AssetStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateAssetStatus(r.Context(), id, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append Outbox Event
	payload, _ := json.Marshal(map[string]interface{}{"asset_id": id, "new_status": req.Status})
	h.repo.AppendEvent(r.Context(), nil, &domain.OutboxEvent{
		Type:    domain.EventAssetTransitioned,
		Payload: payload,
	})

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAsset deletes an asset.
// @Summary Delete Asset
// @Description Permanently removes an asset.
// @Tags Assets
// @Param id path int true "Asset ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventory/assets/{id} [delete]
func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteAsset(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListAssets lists assets filtered by item type.
// @Summary List Assets
// @Description Returns a list of assets, optionally filtered by item_type_id.
// @Tags Assets
// @Produce json
// @Param item_type_id query int true "Item Type ID"
// @Success 200 {array} domain.Asset
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /inventory/assets [get]
func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	itemTypeIDStr := r.URL.Query().Get("item_type_id")
	if itemTypeIDStr == "" {
		http.Error(w, "item_type_id required", http.StatusBadRequest)
		return
	}
	itemTypeID, err := strconv.ParseInt(itemTypeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid item_type_id", http.StatusBadRequest)
		return
	}

	results, err := h.repo.ListAssetsByItemType(r.Context(), itemTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ListRentActions returns all rent actions.
// @Summary List Rent Actions
// @Description Returns all rent actions (reservations).
// @Tags RentActions
// @Produce json
// @Success 200 {array} domain.RentAction
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions [get]
func (h *Handler) ListRentActions(w http.ResponseWriter, r *http.Request) {
	results, err := h.repo.ListRentActions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// RentAction Handlers

// CreateRentAction creates a new reservation request.
// @Summary Create Rent Action
// @Description Creates a new rent action (reservation).
// @Tags RentActions
// @Accept json
// @Produce json
// @Param rent_action body domain.RentAction true "Rent Action Data"
// @Success 201 {object} domain.RentAction
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions [post]
func (h *Handler) CreateRentAction(w http.ResponseWriter, r *http.Request) {
	var ra domain.RentAction
	if err := json.NewDecoder(r.Body).Decode(&ra); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ra.CreatedByUserID = h.getUserIDFromContext(r)

	if err := h.repo.CreateRentAction(r.Context(), &ra); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ra)
}

// GetRentAction retrieves a rent action by ID.
// @Summary Get Rent Action
// @Description Retrieves a rent action by its ID.
// @Tags RentActions
// @Produce json
// @Param id path int true "Rent Action ID"
// @Success 200 {object} domain.RentAction
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions/{id} [get]
func (h *Handler) GetRentAction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/rent-actions/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ra, err := h.repo.GetRentActionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ra == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ra)
}

// SubmitRentAction submits a draft rent action for approval.
// @Summary Submit Rent Action
// @Description Transitions a rent action from Draft to Pending.
// @Tags RentActions
// @Param id path int true "Rent Action ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid State Transition"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions/{id}/submit [post]
func (h *Handler) SubmitRentAction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/rent-actions/")
	idStr = strings.TrimSuffix(idStr, "/submit")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ra, err := h.repo.GetRentActionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ra == nil {
		http.NotFound(w, r)
		return
	}

	if err := ra.Submit(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateRentActionStatus(r.Context(), id, ra.Status, "", time.Time{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApproveRentAction approves a pending rent action.
// @Summary Approve Rent Action
// @Description Transitions a rent action from Pending to Approved. Checks availability.
// @Tags RentActions
// @Param id path int true "Rent Action ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid State Transition"
// @Failure 409 {string} string "Insufficient Inventory"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions/{id}/approve [post]
func (h *Handler) ApproveRentAction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/rent-actions/")
	idStr = strings.TrimSuffix(idStr, "/approve")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ra, err := h.repo.GetRentActionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ra == nil {
		http.NotFound(w, r)
		return
	}

	// Basic Availability Check
	for _, item := range ra.Items {
		if item.ItemKind == "item_type" {
			avail, err := h.repo.GetAvailableQuantity(r.Context(), item.ItemID, ra.StartTime, ra.EndTime)
			if err != nil {
				http.Error(w, "availability check failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if avail < item.RequestedQuantity {
				http.Error(w, fmt.Sprintf("insufficient inventory for item_type %d: requested %d, available %d", item.ItemID, item.RequestedQuantity, avail), http.StatusConflict)
				return
			}
		}
	}

	if err := ra.Approve(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateRentActionStatus(r.Context(), id, ra.Status, "approved_at", *ra.ApprovedAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RejectRentAction rejects a pending rent action.
// @Summary Reject Rent Action
// @Description Transitions a rent action from Pending to Rejected.
// @Tags RentActions
// @Param id path int true "Rent Action ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid State Transition"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions/{id}/reject [post]
func (h *Handler) RejectRentAction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/rent-actions/")
	idStr = strings.TrimSuffix(idStr, "/reject")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ra, err := h.repo.GetRentActionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ra == nil {
		http.NotFound(w, r)
		return
	}

	if err := ra.Reject(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateRentActionStatus(r.Context(), id, ra.Status, "rejected_at", *ra.RejectedAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CancelRentAction cancels a rent action.
// @Summary Cancel Rent Action
// @Description Cancels a rent action.
// @Tags RentActions
// @Param id path int true "Rent Action ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid State Transition"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /rent-actions/{id}/cancel [post]
func (h *Handler) CancelRentAction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/rent-actions/")
	idStr = strings.TrimSuffix(idStr, "/cancel")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ra, err := h.repo.GetRentActionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ra == nil {
		http.NotFound(w, r)
		return
	}

	if err := ra.Cancel(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateRentActionStatus(r.Context(), id, ra.Status, "cancelled_at", *ra.CancelledAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Inspection Handlers

func (h *Handler) CreateInspectionTemplate(w http.ResponseWriter, r *http.Request) {
	var it domain.InspectionTemplate
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateInspectionTemplate(r.Context(), &it); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(it)
}

func (h *Handler) GetRequiredInspections(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/required-inspections")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	a, err := h.repo.GetAssetByID(r.Context(), id)
	if err != nil || a == nil {
		http.NotFound(w, r)
		return
	}

	templates, err := h.repo.GetInspectionTemplatesForItemType(r.Context(), a.ItemTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (h *Handler) SubmitInspection(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/inspections")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var is domain.InspectionSubmission
	if err := json.NewDecoder(r.Body).Decode(&is); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	is.AssetID = assetID

	if err := h.repo.SubmitInspection(r.Context(), &is); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Automatic QC Transition logic (Draft)
	// If a response indicates "QC Passed", we could transition the asset to available.
	// For now, this is a placeholder for future business logic refinement.
	if len(is.Responses) > 0 {
		// Example: if is.Responses[0].Value == "true" { ... }
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(is)
}

// Maintenance Workflow Handlers

func (h *Handler) RecallItemTypeAssets(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/fleet/item-types/")
	idStr = strings.TrimSuffix(idStr, "/recall")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.RecallAssetsByItemType(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RepairAsset(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/repair")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateAssetStatus(r.Context(), id, domain.AssetStatusMaintenance); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefurbishAsset(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/refurbish")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		BuildSpecID int64 `json:"build_spec_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Move to maintenance and assign latest build spec
	a, err := h.repo.GetAssetByID(r.Context(), assetID)
	if err != nil || a == nil {
		http.NotFound(w, r)
		return
	}

	a.Status = domain.AssetStatusMaintenance
	a.CurrentBuildSpecID = &req.BuildSpecID
	a.ProvisioningStatus = domain.ProvisioningFlashing // Transition back to flashing

	if err := h.repo.UpdateAsset(r.Context(), a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Build Spec Handlers

func (h *Handler) CreateBuildSpec(w http.ResponseWriter, r *http.Request) {
	var bs domain.BuildSpec
	if err := json.NewDecoder(r.Body).Decode(&bs); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateBuildSpec(r.Context(), &bs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bs)
}

func (h *Handler) ListBuildSpecs(w http.ResponseWriter, r *http.Request) {
	results, err := h.repo.ListBuildSpecs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Provisioning Handlers

func (h *Handler) StartProvisioning(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/provision")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		BuildSpecID int64  `json:"build_spec_id"`
		PerformedBy string `json:"performed_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pa, err := h.repo.StartProvisioning(r.Context(), assetID, req.BuildSpecID, req.PerformedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pa)
}

func (h *Handler) CompleteProvisioning(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/inventory/assets/")
	idStr = strings.TrimSuffix(idStr, "/complete-provisioning")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_ = assetID

	var req struct {
		ActionID int64  `json:"action_id"`
		Notes    string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAssetRemoteStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/fleet/assets/")
	idStr = strings.TrimSuffix(idStr, "/remote-status")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	asset, err := h.repo.GetAssetByID(r.Context(), id)
	if err != nil || asset == nil {
		http.NotFound(w, r)
		return
	}

	if asset.RemoteManagementID == nil {
		http.Error(w, "asset does not have remote management enabled", http.StatusBadRequest)
		return
	}

	mgr, err := h.remoteRegistry.Get("mock-provider")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := mgr.GetDeviceInfo(r.Context(), *asset.RemoteManagementID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (h *Handler) ApplyAssetRemotePower(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/fleet/assets/")
	idStr = strings.TrimSuffix(idStr, "/remote-power")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Action domain.RemotePowerAction `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	asset, err := h.repo.GetAssetByID(r.Context(), id)
	if err != nil || asset == nil {
		http.NotFound(w, r)
		return
	}

	if asset.RemoteManagementID == nil {
		http.Error(w, "asset does not have remote management enabled", http.StatusBadRequest)
		return
	}

	mgr, err := h.remoteRegistry.Get("mock-provider")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := mgr.ApplyPowerAction(r.Context(), *asset.RemoteManagementID, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
