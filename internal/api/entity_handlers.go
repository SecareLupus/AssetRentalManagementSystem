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

// People & Roles

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var p domain.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreatePerson(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) ListPeople(w http.ResponseWriter, r *http.Request) {
	people, err := h.repo.ListPeople(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/people/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	person, err := h.repo.GetPerson(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if person == nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/people/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var p domain.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.repo.UpdatePerson(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateOrganizationRole(w http.ResponseWriter, r *http.Request) {
	var or domain.OrganizationRole
	if err := json.NewDecoder(r.Body).Decode(&or); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateOrganizationRole(r.Context(), &or); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(or)
}

func (h *Handler) ListOrganizationRoles(w http.ResponseWriter, r *http.Request) {
	var orgID *int64
	if cidStr := r.URL.Query().Get("organization_id"); cidStr != "" {
		if val, err := strconv.ParseInt(cidStr, 10, 64); err == nil {
			orgID = &val
		}
	}
	var personID *int64
	if pidStr := r.URL.Query().Get("person_id"); pidStr != "" {
		if val, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			personID = &val
		}
	}
	roles, err := h.repo.ListOrganizationRoles(r.Context(), orgID, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(roles)
}

// Places

func (h *Handler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	var p domain.Place
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreatePlace(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) ListPlaces(w http.ResponseWriter, r *http.Request) {
	var ownerID *int64
	if oidStr := r.URL.Query().Get("owner_id"); oidStr != "" {
		if val, err := strconv.ParseInt(oidStr, 10, 64); err == nil {
			ownerID = &val
		}
	}
	var parentID *int64
	if pidStr := r.URL.Query().Get("parent_id"); pidStr != "" {
		if val, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			parentID = &val
		}
	}
	places, err := h.repo.ListPlaces(r.Context(), ownerID, parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(places)
}

func (h *Handler) GetPlace(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/places/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	place, err := h.repo.GetPlace(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if place == nil {
		http.Error(w, "place not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(place)
}

func (h *Handler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/places/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var p domain.Place
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.repo.UpdatePlace(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/events/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	event, err := h.repo.GetEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if event == nil {
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(event)
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

// Demands (Successor to EventAssetNeeds)

func (h *Handler) ListEventDemands(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/events/")
	idStr = strings.TrimSuffix(idStr, "/demands")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}
	demands, err := h.repo.ListDemandsByEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(demands)
}

func (h *Handler) UpdateEventDemands(w http.ResponseWriter, r *http.Request) {
	var demands []domain.Demand
	if err := json.NewDecoder(r.Body).Decode(&demands); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, d := range demands {
		if d.ID > 0 {
			h.repo.UpdateDemand(r.Context(), &d)
		} else {
			h.repo.CreateDemand(r.Context(), &d)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete handlers

func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/companies/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteCompany(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/people/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeletePerson(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/places/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeletePlace(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/entities/events/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteEvent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
