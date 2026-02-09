package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type IngestTargetModel string

const (
	IngestTargetItemType IngestTargetModel = "item_type"
	IngestTargetAsset    IngestTargetModel = "asset"
	IngestTargetCompany  IngestTargetModel = "company"
	IngestTargetPerson   IngestTargetModel = "person"
	IngestTargetPlace    IngestTargetModel = "place"
)

type IngestAuthType string

const (
	IngestAuthNone   IngestAuthType = "none"
	IngestAuthBearer IngestAuthType = "bearer"
)

type IngestSource struct {
	ID                  int64           `json:"id"`
	Name                string          `json:"name"`
	BaseURL             string          `json:"base_url"`
	AuthType            IngestAuthType  `json:"auth_type"`
	AuthEndpoint        string          `json:"auth_endpoint"`
	VerifyEndpoint      string          `json:"verify_endpoint,omitempty"`
	RefreshEndpoint     string          `json:"refresh_endpoint"`
	AuthCredentials     json.RawMessage `json:"auth_credentials,omitempty"`
	LastToken           string          `json:"last_token"`
	RefreshToken        string          `json:"refresh_token"`
	TokenExpiry         *time.Time      `json:"token_expiry"`
	SyncIntervalSeconds int             `json:"sync_interval_seconds"`
	IsActive            bool            `json:"is_active"`

	LastStatus string `json:"last_status"`
	LastError  string `json:"last_error"`

	Endpoints []IngestEndpoint `json:"endpoints,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IngestEndpoint struct {
	ID           int64           `json:"id"`
	SourceID     int64           `json:"source_id"`
	Path         string          `json:"path"`
	Method       string          `json:"method"`
	RequestBody  json.RawMessage `json:"request_body,omitempty"`
	RespStrategy string          `json:"resp_strategy"` // 'single', 'list', 'auto'
	ItemsPath    string          `json:"items_path"`    // JSONPath to the list of items
	IsActive     bool            `json:"is_active"`

	LastSyncAt      *time.Time `json:"last_sync_at"`
	LastSuccessAt   *time.Time `json:"last_success_at"`
	LastETag        string     `json:"last_etag"`
	LastPayloadHash string     `json:"last_payload_hash"`

	Mappings []IngestMapping `json:"mappings,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InferredField struct {
	Path             string `json:"path"`
	Label            string `json:"label"`
	Type             string `json:"type"`
	SuggestedMapping string `json:"suggest_mapping,omitempty"`
	SuggestedModel   string `json:"suggested_model,omitempty"`
	IsIdentity       bool   `json:"is_identity"`
}

type DiscoveryResponse struct {
	RawResponse    json.RawMessage `json:"raw_response"`
	SampleItems    []interface{}   `json:"sample_items,omitempty"`
	ItemsPath      string          `json:"items_path,omitempty"`
	InferredFields []InferredField `json:"inferred_fields,omitempty"`
}

type IngestMapping struct {
	ID          int64             `json:"id"`
	EndpointID  int64             `json:"endpoint_id"`
	JSONPath    string            `json:"json_path"`    // e.g. "$.sku"
	TargetModel IngestTargetModel `json:"target_model"` // e.g. "asset"
	TargetField string            `json:"target_field"` // e.g. "code"
	IsIdentity  bool              `json:"is_identity"`  // For UPSERT
	CreatedAt   time.Time         `json:"created_at"`
}

// UnwrapJSON checks if a json.RawMessage contains a double-encoded JSON string
// and returns the inner JSON if valid.
func UnwrapJSON(m json.RawMessage) json.RawMessage {
	if len(m) == 0 {
		return m
	}
	var s string
	if err := json.Unmarshal(m, &s); err == nil {
		if json.Valid([]byte(s)) {
			return json.RawMessage(s)
		}
	}
	return m
}

// DiscoverTokens attempts to find access and refresh tokens in a generic map
// using common field name variations.
func DiscoverTokens(data map[string]interface{}) (accessToken, refreshToken string, expiresIn int) {
	// Access Token variations
	accessKeys := []string{"access_token", "api_token", "token", "accessToken", "jwt"}
	for _, k := range accessKeys {
		if v, ok := data[k].(string); ok && v != "" {
			accessToken = v
			break
		}
	}

	// Refresh Token variations
	refreshKeys := []string{"refresh_token", "refreshToken", "refresh"}
	for _, k := range refreshKeys {
		if v, ok := data[k].(string); ok && v != "" {
			refreshToken = v
			break
		}
	}

	// Expires In variations
	expiryKeys := []string{"expires_in", "expiresIn", "expires"}
	for _, k := range expiryKeys {
		if v, ok := data[k].(float64); ok {
			expiresIn = int(v)
			break
		}
		if v, ok := data[k].(int); ok {
			expiresIn = v
			break
		}
	}

	return
}

// DiscoverSchema attempts to find a list of items and infer their fields.
func DiscoverSchema(body []byte) (*DiscoveryResponse, error) {
	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	resp := &DiscoveryResponse{
		RawResponse: json.RawMessage(body),
	}

	var items []interface{}
	itemsPath := "$"

	// Heuristic 1: Is it already a list?
	if list, ok := raw.([]interface{}); ok {
		items = list
	} else if m, ok := raw.(map[string]interface{}); ok {
		// Heuristic 2: Look for common list keys
		listKeys := []string{"data", "items", "results", "records", "devices"}
		for _, k := range listKeys {
			if v, ok := m[k].([]interface{}); ok {
				items = v
				itemsPath = "$." + k
				break
			}
		}
	}

	if len(items) == 0 {
		return resp, nil // No items found to infer from
	}

	resp.ItemsPath = itemsPath
	resp.SampleItems = items
	if len(items) > 5 {
		resp.SampleItems = items[:5]
	}

	// Infer fields from the first item
	if first, ok := items[0].(map[string]interface{}); ok {
		for k, v := range first {
			field := InferredField{
				Path:  "$." + k,
				Label: strings.Title(strings.ReplaceAll(k, "_", " ")),
			}

			switch v.(type) {
			case string:
				field.Type = "string"
				// Basic mapping heuristics
				lowK := strings.ToLower(k)
				if strings.Contains(lowK, "id") || strings.Contains(lowK, "tag") || strings.Contains(lowK, "token") || strings.Contains(lowK, "serial") {
					field.IsIdentity = true
				}
				if strings.Contains(lowK, "name") {
					field.SuggestedMapping = "name"
					field.SuggestedModel = "asset"
				} else if strings.Contains(lowK, "tag") {
					field.SuggestedMapping = "asset_tag"
					field.SuggestedModel = "asset"
				} else if strings.Contains(lowK, "serial") {
					field.SuggestedMapping = "serial_number"
					field.SuggestedModel = "asset"
				}
			case float64:
				field.Type = "number"
			case bool:
				field.Type = "boolean"
			case map[string]interface{}:
				field.Type = "object"
			case []interface{}:
				field.Type = "array"
			default:
				field.Type = "unknown"
			}
			resp.InferredFields = append(resp.InferredFields, field)
		}
	}

	return resp, nil
}
