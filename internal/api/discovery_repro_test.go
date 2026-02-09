package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Discovery_403Propagation(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	// Mock source and endpoint
	src := &domain.IngestSource{
		ID:       1,
		BaseURL:  "", // Will be set to test server URL
		AuthType: domain.IngestAuthNone,
	}
	ep := &domain.IngestEndpoint{
		ID:       1,
		SourceID: 1,
		Path:     "/devices/list",
		Method:   "GET",
	}

	// Mock upstream server returning 403
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": 403,
			"message":    "Forbidden resource",
		})
	}))
	defer upstream.Close()

	src.BaseURL = upstream.URL

	repo.On("GetIngestEndpoint", mock.Anything, int64(1)).Return(ep, nil)
	repo.On("GetIngestSource", mock.Anything, int64(1)).Return(src, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/admin/ingest/endpoints/1/discovery", nil)
	w := httptest.NewRecorder()

	h.Discovery(w, req)

	// CURRENT BEHAVIOR: Returns 200 OK because io.Copy is used without checking status
	// DESIRED BEHAVIOR: Returns 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code, "Should propagate 403 status code")

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, float64(403), body["statusCode"])
}

func TestHandler_Discovery_403Retry(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	// Mock source and endpoint
	src := &domain.IngestSource{
		ID:        1,
		BaseURL:   "",
		AuthType:  domain.IngestAuthBearer,
		LastToken: "old-token",
	}
	ep := &domain.IngestEndpoint{
		ID:       1,
		SourceID: 1,
		Path:     "/devices/list",
		Method:   "GET",
	}

	callCount := 0
	// Mock upstream server: first call 403, second call 200 after refresh
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			assert.Equal(t, "Bearer old-token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		assert.Equal(t, "Bearer new-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
	}))
	defer upstream.Close()

	src.BaseURL = upstream.URL
	src.AuthEndpoint = upstream.URL + "/auth"

	// Mock auth server response
	// Note: authenticateSource uses POST to authURL
	// We handle this in the same upstream for simplicity
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "new-token",
		})
	}))
	defer authServer.Close()
	src.AuthEndpoint = authServer.URL

	repo.On("GetIngestEndpoint", mock.Anything, int64(1)).Return(ep, nil)
	repo.On("GetIngestSource", mock.Anything, int64(1)).Return(src, nil)
	repo.On("UpdateIngestSource", mock.Anything, mock.MatchedBy(func(s *domain.IngestSource) bool {
		return s.LastToken == "new-token"
	})).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/admin/ingest/endpoints/1/discovery", nil)
	w := httptest.NewRecorder()

	h.Discovery(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, callCount, "Should have retried the request")
	repo.AssertExpectations(t)
}

func TestHandler_Discovery_NullBodyGET(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	// Mock source and endpoint
	src := &domain.IngestSource{
		ID:       1,
		BaseURL:  "",
		AuthType: domain.IngestAuthNone,
	}
	ep := &domain.IngestEndpoint{
		ID:          1,
		SourceID:    1,
		Path:        "/devices/list",
		Method:      "GET",
		RequestBody: json.RawMessage(`null`),
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		body, _ := io.ReadAll(r.Body)
		assert.Empty(t, body, "Body should be empty for GET with null RequestBody")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
	}))
	defer upstream.Close()

	src.BaseURL = upstream.URL

	repo.On("GetIngestEndpoint", mock.Anything, int64(1)).Return(ep, nil)
	repo.On("GetIngestSource", mock.Anything, int64(1)).Return(src, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/admin/ingest/endpoints/1/discovery", nil)
	w := httptest.NewRecorder()

	h.Discovery(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}
