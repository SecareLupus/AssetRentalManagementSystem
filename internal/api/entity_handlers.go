package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/desmond/rental-management-system/internal/domain"
)

// Companies

func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var c domain.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateCompany(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(companies)
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/companies/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	company, err := h.repo.GetCompany(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if company == nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(company)
}

func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/companies/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var c domain.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	if err := h.repo.UpdateCompany(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Contacts

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var c domain.Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateContact(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request) {
	var companyID *int64
	if cidStr := r.URL.Query().Get("company_id"); cidStr != "" {
		if val, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			companyID = &val
		}
	}
	contacts, err := h.repo.ListContacts(r.Context(), companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(contacts)
}

// Sites

func (h *Handler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var s domain.Site
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateSite(r.Context(), &s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) ListSites(w http.ResponseWriter, r *http.Request) {
	var companyID *int64
	if cidStr := r.URL.Query().Get("company_id"); cidStr != "" {
		if val, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			companyID = &val
		}
	}
	sites, err := h.repo.ListSites(r.Context(), companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sites)
}

// Locations

func (h *Handler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var l domain.Location
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateLocation(r.Context(), &l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(l)
}

func (h *Handler) ListLocations(w http.ResponseWriter, r *http.Request) {
	var siteID *int64
	if sidStr := r.URL.Query().Get("site_id"); sidStr != "" {
		if val, err := strconv.ParseInt(sidStr, 10, 64); err == nil {
			siteID = &val
		}
	}
	var parentID *int64
	if pidStr := r.URL.Query().Get("parent_id"); pidStr != "" {
		if val, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			parentID = &val
		}
	}
	locations, err := h.repo.ListLocations(r.Context(), siteID, parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(locations)
}

// Events

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateEvent(r.Context(), &e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	var companyID *int64
	if cidStr := r.URL.Query().Get("company_id"); cidStr != "" {
		if val, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			companyID = &val
		}
	}
	events, err := h.repo.ListEvents(r.Context(), companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/events/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e.ID = id
	if err := h.repo.UpdateEvent(r.Context(), &e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// EventAssetNeeds

func (h *Handler) ListEventAssetNeeds(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/events/")
	idStr = strings.TrimSuffix(idStr, "/needs")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}
	needs, err := h.repo.ListEventAssetNeeds(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(needs)
}

func (h *Handler) UpdateEventAssetNeeds(w http.ResponseWriter, r *http.Request) {
	var needs []domain.EventAssetNeed
	if err := json.NewDecoder(r.Body).Decode(&needs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, n := range needs {
		if n.ID > 0 {
			h.repo.UpdateEventAssetNeed(r.Context(), &n)
		} else {
			h.repo.CreateEventAssetNeed(r.Context(), &n)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
