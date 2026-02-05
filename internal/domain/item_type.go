package domain

import (
	"encoding/json"
	"time"
)

type ItemKind string

const (
	ItemKindSerialized ItemKind = "serialized"
	ItemKindFungible   ItemKind = "fungible"
	ItemKindKit        ItemKind = "kit"
)

type LifecycleFeatures struct {
	RemoteManagement  bool `json:"remote_management"`
	Provisioning      bool `json:"provisioning"`
	Refurbishment     bool `json:"refurbishment"`
	BuildSpecTracking bool `json:"build_spec_tracking"`
}

// ItemType represents a template for a piece of equipment.
// It maps to https://schema.org/ProductModel
type ItemType struct {
	ID                int64             `json:"id"`
	Code              string            `json:"code"` // Maps to schema.org/sku
	Name              string            `json:"name"` // Maps to schema.org/name
	Kind              ItemKind          `json:"kind"` // Maps to schema.org/category
	IsActive          bool              `json:"is_active"`
	SupportedFeatures LifecycleFeatures `json:"supported_features"`
	CreatedByUserID   *int64            `json:"created_by_user_id,omitempty"`
	UpdatedByUserID   *int64            `json:"updated_by_user_id,omitempty"`
	SchemaOrg         json.RawMessage   `json:"schema_org,omitempty" swaggertype:"string" example:"{}"`
	Metadata          json.RawMessage   `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
