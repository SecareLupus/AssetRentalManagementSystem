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

// ItemType represents a template for a piece of equipment.
// It maps to https://schema.org/ProductModel
type ItemType struct {
	ID        int64           `json:"id"`
	Code      string          `json:"code"` // Maps to schema.org/sku
	Name      string          `json:"name"` // Maps to schema.org/name
	Kind      ItemKind        `json:"kind"` // Maps to schema.org/category
	IsActive  bool            `json:"is_active"`
	SchemaOrg json.RawMessage `json:"schema_org,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
