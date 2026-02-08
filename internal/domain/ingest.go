package domain

import (
	"encoding/json"
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
	IsActive     bool            `json:"is_active"`

	LastSyncAt      *time.Time `json:"last_sync_at"`
	LastSuccessAt   *time.Time `json:"last_success_at"`
	LastETag        string     `json:"last_etag"`
	LastPayloadHash string     `json:"last_payload_hash"`

	Mappings []IngestMapping `json:"mappings,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
