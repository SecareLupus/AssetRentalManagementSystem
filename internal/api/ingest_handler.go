package api

import (
	"bytes"
	"context"
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

	if err := h.authenticateSource(r.Context(), src); err != nil {
		http.Error(w, "auth failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"token":  src.LastToken,
	})
}

func (h *Handler) authenticateSource(ctx context.Context, src *domain.IngestSource) error {
	src.AuthCredentials = domain.UnwrapJSON(src.AuthCredentials)

	// Perform auth request
	authURL := src.AuthEndpoint
	if !strings.HasPrefix(authURL, "http") {
		authURL = src.BaseURL + authURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(authURL, "application/json", bytes.NewReader(src.AuthCredentials))
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed with %s: %s", resp.Status, body)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode auth response")
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
			return fmt.Errorf("failed to create verify request")
		}
		req.Header.Set("Content-Type", "application/json")
		if loginToken != "" {
			req.Header.Set("Authorization", "Bearer "+loginToken)
		}

		vResp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("verify request failed: %w", err)
		}
		defer vResp.Body.Close()

		if vResp.StatusCode < 200 || vResp.StatusCode >= 300 {
			vBody, _ := io.ReadAll(vResp.Body)
			return fmt.Errorf("verify failed with %s: %s", vResp.Status, vBody)
		}

		if err := json.NewDecoder(vResp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode verify response")
		}
	}

	accessToken, refreshToken, expiresIn := domain.DiscoverTokens(data)

	if accessToken == "" {
		return fmt.Errorf("no access token found in response")
	}

	src.LastToken = accessToken
	if refreshToken != "" {
		src.RefreshToken = refreshToken
	}
	if expiresIn > 0 {
		expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)
		src.TokenExpiry = &expiry
	}
	return h.repo.UpdateIngestSource(ctx, src)
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

	// Refined body handling: Don't send "null" or empty body for GET/DELETE,
	// or if the body is explicitly the JSON literal "null".
	var bodyReader io.Reader
	sendBody := false
	if len(body) > 0 && string(body) != "null" {
		method := strings.ToUpper(targetEP.Method)
		if method == "POST" || method == "PUT" || method == "PATCH" {
			sendBody = true
			bodyReader = bytes.NewReader(body)
		}
	}

	apiReq, _ := http.NewRequestWithContext(r.Context(), targetEP.Method, fullURL, bodyReader)
	apiReq.Header.Set("Accept", "application/json")
	if sendBody {
		apiReq.Header.Set("Content-Type", "application/json")
	}

	if targetSrc.AuthType == domain.IngestAuthBearer {
		// Check if token needs refresh
		if targetSrc.LastToken == "" || (targetSrc.TokenExpiry != nil && time.Now().Add(5*time.Minute).After(*targetSrc.TokenExpiry)) {
			// Refresh token
			if err := h.authenticateSource(r.Context(), targetSrc); err != nil {
				// Log error but try anyway? Or fail?
				// Better to fail because we know it's expired
				http.Error(w, "failed to refresh token: "+err.Error(), http.StatusUnauthorized)
				return
			}
		}
		apiReq.Header.Set("Authorization", "Bearer "+targetSrc.LastToken)
	}

	resp, err := client.Do(apiReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Retry on 401 or 403 if using Bearer auth
	if (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) && targetSrc.AuthType == domain.IngestAuthBearer {
		// Consuming the body is good practice before closing, though we are about to close it
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// Attempt refresh
		if err := h.authenticateSource(r.Context(), targetSrc); err != nil {
			// If refresh fails, return the original error or the refresh error
			http.Error(w, "upstream 401 and refresh failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Re-create request with new token
		retryReq, _ := http.NewRequestWithContext(r.Context(), targetEP.Method, fullURL, bodyReader)
		retryReq.Header.Set("Accept", "application/json")
		if sendBody {
			retryReq.Header.Set("Content-Type", "application/json")
		}
		retryReq.Header.Set("Authorization", "Bearer "+targetSrc.LastToken)

		resp, err = client.Do(retryReq)
		if err != nil {
			http.Error(w, "retry failed: "+err.Error(), http.StatusBadGateway)
			return
		}
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	disco, err := domain.DiscoverSchema(respBody)
	if err != nil {
		// Fallback to raw response if discovery fails
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(disco)
}

func (h *Handler) SyncSourceNow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "sync triggered"})
}
