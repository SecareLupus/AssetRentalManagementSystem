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
)

type Handler struct {
	repo db.Repository
}

func NewHandler(repo db.Repository) *Handler {
	return &Handler{repo: repo}
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

	if err := h.repo.CreateAsset(r.Context(), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

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

	w.WriteHeader(http.StatusNoContent)
}

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

// RentAction Handlers

func (h *Handler) CreateRentAction(w http.ResponseWriter, r *http.Request) {
	var ra domain.RentAction
	if err := json.NewDecoder(r.Body).Decode(&ra); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateRentAction(r.Context(), &ra); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ra)
}

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

	if err := h.repo.CompleteProvisioning(r.Context(), req.ActionID, req.Notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
