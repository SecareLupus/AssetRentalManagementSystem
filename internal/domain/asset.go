package domain

import (
	"encoding/json"
	"time"
)

type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "available"
	AssetStatusReserved    AssetStatus = "reserved"
	AssetStatusMaintenance AssetStatus = "maintenance"
	AssetStatusRetired     AssetStatus = "retired"
)

type Asset struct {
	ID                int64           `json:"id"`
	ItemTypeID        int64           `json:"item_type_id"`
	AssetTag          *string         `json:"asset_tag,omitempty"`
	SerialNumber      *string         `json:"serial_number,omitempty"`
	Status            AssetStatus     `json:"status"`
	Location          *string         `json:"location,omitempty"`
	AssignedTo        *string         `json:"assigned_to,omitempty"`
	MeshNodeID        *string         `json:"mesh_node_id,omitempty"`
	WireguardHostname *string         `json:"wireguard_hostname,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
