package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	src.AuthCredentials = domain.UnwrapJSON(src.AuthCredentials)

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

	src.AuthCredentials = domain.UnwrapJSON(src.AuthCredentials)

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

// Endpoints

func (h *Handler) CreateIngestEndpoint(w http.ResponseWriter, r *http.Request) {
	var ep domain.IngestEndpoint
	if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ep.RequestBody = domain.UnwrapJSON(ep.RequestBody)

	if err := h.repo.CreateIngestEndpoint(r.Context(), &ep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ep)
}

func (h *Handler) UpdateIngestEndpoint(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/endpoints/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ep domain.IngestEndpoint
	if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	ep.ID = id

	ep.RequestBody = domain.UnwrapJSON(ep.RequestBody)

	if err := h.repo.UpdateIngestEndpoint(r.Context(), &ep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteIngestEndpoint(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/endpoints/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteIngestEndpoint(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetEndpointMappings(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/endpoints/")
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

	if err := h.repo.SetEndpointMappings(r.Context(), id, mappings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) TestAuth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SourceID int64 `json:"source_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.SourceID == 0 {
		http.Error(w, "source_id is required", http.StatusBadRequest)
		return
	}

	src, err := h.repo.GetIngestSource(r.Context(), req.SourceID)
	if err != nil {
		http.Error(w, "repository error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if src == nil {
		http.Error(w, "source not found", http.StatusNotFound)
		return
	}

	if src.AuthType != domain.IngestAuthBearer {
		json.NewEncoder(w).Encode(map[string]string{"status": "no auth needed"})
		return
	}

	src.AuthCredentials = domain.UnwrapJSON(src.AuthCredentials)

	// Perform auth request
	authURL := src.AuthEndpoint
	if !strings.HasPrefix(authURL, "http") {
		authURL = src.BaseURL + authURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(authURL, "application/json", bytes.NewReader(src.AuthCredentials))
	if err != nil {
		http.Error(w, "auth request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("auth failed with %s: %s", resp.Status, body), http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "failed to decode auth response", http.StatusInternalServerError)
		return
	}

	// Step 1.5: Verify (Optional Triple-Auth)
	if src.VerifyEndpoint != "" {
		verifyURL := src.VerifyEndpoint
		if !strings.HasPrefix(verifyURL, "http") {
			verifyURL = src.BaseURL + verifyURL
		}

		// Attempt to find a token to use for verification if needed
		loginToken, _, _ := domain.DiscoverTokens(data)

		verifyBody, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", verifyURL, bytes.NewReader(verifyBody))
		if err != nil {
			http.Error(w, "failed to create verify request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if loginToken != "" {
			req.Header.Set("Authorization", "Bearer "+loginToken)
		}

		vResp, err := client.Do(req)
		if err != nil {
			http.Error(w, "verify request failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer vResp.Body.Close()

		if vResp.StatusCode < 200 || vResp.StatusCode >= 300 {
			vBody, _ := io.ReadAll(vResp.Body)
			http.Error(w, fmt.Sprintf("verify failed with %s: %s", vResp.Status, vBody), http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(vResp.Body).Decode(&data); err != nil {
			http.Error(w, "failed to decode verify response", http.StatusInternalServerError)
			return
		}
	}

	accessToken, refreshToken, expiresIn := domain.DiscoverTokens(data)

	if accessToken == "" {
		http.Error(w, "no access token found in response", http.StatusBadGateway)
		return
	}

	src.LastToken = accessToken
	if refreshToken != "" {
		src.RefreshToken = refreshToken
	}
	if expiresIn > 0 {
		expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)
		src.TokenExpiry = &expiry
	}
	h.repo.UpdateIngestSource(r.Context(), src)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"token":  src.LastToken,
	})
}

func (h *Handler) Discovery(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/ingest/endpoints/")
	idStr = strings.TrimSuffix(idStr, "/discovery")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	targetEP, err := h.repo.GetIngestEndpoint(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if targetEP == nil {
		http.NotFound(w, r)
		return
	}

	targetSrc, err := h.repo.GetIngestSource(r.Context(), targetEP.SourceID)
	if err != nil || targetSrc == nil {
		http.Error(w, "parent source not found", http.StatusNotFound)
		return
	}

	fullURL := targetSrc.BaseURL + targetEP.Path
	client := &http.Client{Timeout: 30 * time.Second}

	body := domain.UnwrapJSON(targetEP.RequestBody)
	apiReq, _ := http.NewRequestWithContext(r.Context(), targetEP.Method, fullURL, bytes.NewReader(body))

	if targetSrc.AuthType == domain.IngestAuthBearer && targetSrc.LastToken != "" {
		apiReq.Header.Set("Authorization", "Bearer "+targetSrc.LastToken)
	}

	resp, err := client.Do(apiReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (h *Handler) SyncSourceNow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "sync triggered"})
}
