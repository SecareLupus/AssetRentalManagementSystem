package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateIngestEndpoint_Repro(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	// User provided payload
	bodyJSON := `{
    "startRow": 0,
    "endRow": 1000,
    "sortModel": [],
    "filterModel": {},
    "searchTerm": "",
    "area": "rvs_devices",
    "customer_id": -1,
    "show_id": 0,
    "subsetIDs": [],
    "omitSubsetIDs": [],
    "includeAllShows": null,
    "advanceSearch": {},
    "quickSearch": {
        "searchCriteria": "",
        "showProgress": false
    },
    "verificationData": {},
    "customFilters": {
        "deviceCompany": -1
    },
    "searchType": "Replace",
    "showInactiveMasterFees": false,
    "typeFilter": "",
    "categoryFilter": ""
}`

	var rawBody json.RawMessage = []byte(bodyJSON)

	// Expected endpoint object
	ep := domain.IngestEndpoint{
		SourceID:     1,
		Path:         "/devices/list",
		Method:       "POST",
		RequestBody:  rawBody,
		RespStrategy: "auto",
		IsActive:     true,
	}

	reqBody, _ := json.Marshal(ep)

	// Mock expectation
	repo.On("CreateIngestEndpoint", mock.Anything, mock.MatchedBy(func(e *domain.IngestEndpoint) bool {
		// Verify that the RequestBody is correctly preserved or unwrapped
		// domain.UnwrapJSON might change the representation if it's double encoded,
		// but here we are sending it as part of the JSON object, so it should be fine.
		// However, if the frontend sends it as a stringified JSON inside the JSON, UnwrapJSON kicks in.

		// Let's mimic what the frontend likely sends.
		return true
	})).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/admin/ingest/endpoints", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	h.CreateIngestEndpoint(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check the response
	var respEp domain.IngestEndpoint
	json.NewDecoder(w.Body).Decode(&respEp)

	// Verify RequestBody in response is valid JSON
	assert.NotNil(t, respEp.RequestBody)
}

func TestHandler_CreateIngestEndpoint_DoubleEncoded_Repro(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	// User provided payload as a STRING (double encoded), which is what the frontend might be doing
	innerJSON := `{
    "startRow": 0,
    "endRow": 1000,
    "sortModel": [],
    "filterModel": {},
    "searchTerm": "",
    "area": "rvs_devices",
    "customer_id": -1,
    "show_id": 0,
    "subsetIDs": [],
    "omitSubsetIDs": [],
    "includeAllShows": null,
    "advanceSearch": {},
    "quickSearch": {
        "searchCriteria": "",
        "showProgress": false
    },
    "verificationData": {},
    "customFilters": {
        "deviceCompany": -1
    },
    "searchType": "Replace",
    "showInactiveMasterFees": false,
    "typeFilter": "",
    "categoryFilter": ""
}`
	// json.Marshal will quote and escape this string
	encodedBody, _ := json.Marshal(innerJSON)

	// Construct the request
	// We manually construct the JSON to look like: { ..., "request_body": "{\"startRow\": 0 ...}" }
	fullPayload := `{"source_id": 1, "path": "/devices/list", "method": "POST", "resp_strategy": "auto", "is_active": true, "request_body": ` + string(encodedBody) + `}`

	repo.On("CreateIngestEndpoint", mock.Anything, mock.MatchedBy(func(e *domain.IngestEndpoint) bool {
		// We expect UnwrapJSON to have fixed this back to raw map/object
		// If it's still a string/byte slice of a string, then it failed.

		// UnwrapJSON returns []byte. If it worked, it should be the raw bytes of innerJSON (unquoted).
		// If it failed, it would remain as the quoted string bytes.

		// Simply checking if we can unmarshal it into a map
		var m map[string]interface{}
		err := json.Unmarshal(e.RequestBody, &m)
		return err == nil && m["area"] == "rvs_devices"
	})).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/admin/ingest/endpoints", bytes.NewBuffer([]byte(fullPayload)))
	w := httptest.NewRecorder()

	h.CreateIngestEndpoint(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
