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
	ID                  int64             `json:"id"`
	Name                string            `json:"name"`
	TargetModel         IngestTargetModel `json:"target_model"`
	APIURL              string            `json:"api_url"`
	AuthType            IngestAuthType    `json:"auth_type"`
	AuthCredentials     json.RawMessage   `json:"auth_credentials,omitempty"`
	SyncIntervalSeconds int               `json:"sync_interval_seconds"`
	IsActive            bool              `json:"is_active"`

	LastSyncAt    *time.Time `json:"last_sync_at"`
	LastSuccessAt *time.Time `json:"last_success_at"`
	LastStatus    string     `json:"last_status"`
	LastError     string     `json:"last_error"`
	NextSyncAt    *time.Time `json:"next_sync_at"`

	LastETag        string `json:"last_etag"`
	LastPayloadHash string `json:"last_payload_hash"`

	Mappings []IngestMapping `json:"mappings,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IngestMapping struct {
	ID          int64     `json:"id"`
	SourceID    int64     `json:"source_id"`
	JSONPath    string    `json:"json_path"`    // e.g. "$.sku"
	TargetField string    `json:"target_field"` // e.g. "code"
	IsIdentity  bool      `json:"is_identity"`  // For UPSERT
	CreatedAt   time.Time `json:"created_at"`
}
