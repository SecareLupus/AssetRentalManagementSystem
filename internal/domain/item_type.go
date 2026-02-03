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

type ItemType struct {
	ID        int64           `json:"id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	Kind      ItemKind        `json:"kind"`
	IsActive  bool            `json:"is_active"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
