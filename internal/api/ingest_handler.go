package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
)

func (h *Handler) ListIngestSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.repo.ListIngestSources(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

func (h *Handler) CreateIngestSource(w http.ResponseWriter, r *http.Request) {
	var src domain.IngestSource
	if err := json.NewDecoder(r.Body).Decode(&src); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateIngestSource(r.Context(), &src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(src)
}

func (h *Handler) GetIngestSource(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/sources/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	src, err := h.repo.GetIngestSource(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if src == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(src)
}

func (h *Handler) UpdateIngestSource(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/sources/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var src domain.IngestSource
	if err := json.NewDecoder(r.Body).Decode(&src); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	src.ID = id

	if err := h.repo.UpdateIngestSource(r.Context(), &src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteIngestSource(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/sources/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteIngestSource(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetIngestMappings(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/sources/")
	idStr = strings.TrimSuffix(idStr, "/mappings")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var mappings []domain.IngestMapping
	if err := json.NewDecoder(r.Body).Decode(&mappings); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.SetIngestMappings(r.Context(), id, mappings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PreviewSource fetches sample data from an endpoint to allow schema discovery.
func (h *Handler) PreviewSource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		APIURL          string                `json:"api_url"`
		AuthType        domain.IngestAuthType `json:"auth_type"`
		AuthCredentials string                `json:"auth_credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}

	apiReq, err := http.NewRequestWithContext(r.Context(), "GET", req.APIURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuthType == domain.IngestAuthBearer && req.AuthCredentials != "" {
		apiReq.Header.Set("Authorization", "Bearer "+req.AuthCredentials)
	}

	resp, err := client.Do(apiReq)
	if err != nil {
		http.Error(w, "failed to reach endpoint: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, "endpoint returned "+resp.Status+": "+string(body), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (h *Handler) SyncSourceNow(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/sources/")
	idStr = strings.TrimSuffix(idStr, "/sync")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	src, err := h.repo.GetIngestSource(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if src == nil {
		http.NotFound(w, r)
		return
	}

	// For now, we'll just acknowledge the request.
	// In a real implementation, this would trigger the worker immediately.
	// Since the worker isn't implemented yet, we'll return a 202 Accepted.
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "sync triggered"})
}
