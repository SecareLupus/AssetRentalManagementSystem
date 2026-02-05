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
	AssetStatusDeployed    AssetStatus = "deployed"
	AssetStatusRecalled    AssetStatus = "recalled"
)

type ProvisioningStatus string

const (
	ProvisioningUnprovisioned ProvisioningStatus = "unprovisioned"
	ProvisioningFlashing      ProvisioningStatus = "flashing"
	ProvisioningConfigured    ProvisioningStatus = "configured"
	ProvisioningReady         ProvisioningStatus = "ready"
)

// Asset represents a specific physical item.
// It maps to https://schema.org/IndividualProduct
type Asset struct {
	ID                int64       `json:"id"`
	ItemTypeID        int64       `json:"item_type_id"`            // Links to the ProductModel (ItemType)
	AssetTag          *string     `json:"asset_tag,omitempty"`     // Maps to schema.org/identifier
	SerialNumber      *string     `json:"serial_number,omitempty"` // Maps to schema.org/serialNumber
	Status            AssetStatus `json:"status"`
	Location          *string     `json:"location,omitempty"`
	AssignedTo        *string     `json:"assigned_to,omitempty"`
	MeshNodeID        *string     `json:"mesh_node_id,omitempty"`
	WireguardHostname *string     `json:"wireguard_hostname,omitempty"`
	ManagementURL     *string     `json:"management_url,omitempty"` // Link to external management tool

	// Fleet Management Fields
	BuildSpecVersion   *string            `json:"build_spec_version,omitempty"`
	ProvisioningStatus ProvisioningStatus `json:"provisioning_status,omitempty"`
	FirmwareVersion    *string            `json:"firmware_version,omitempty"`
	Hostname           *string            `json:"hostname,omitempty"`
	RemoteManagementID *string            `json:"remote_management_id,omitempty"`
	CurrentBuildSpecID *int64             `json:"current_build_spec_id,omitempty"`
	LastInspectionAt   *time.Time         `json:"last_inspection_at,omitempty"`
	UsageHours         float64            `json:"usage_hours"`
	NextServiceHours   float64            `json:"next_service_hours"`
	CreatedByUserID    *int64             `json:"created_by_user_id,omitempty"` // Audit trail
	UpdatedByUserID    *int64             `json:"updated_by_user_id,omitempty"` // Audit trail

	SchemaOrg json.RawMessage `json:"schema_org,omitempty" swaggertype:"string" example:"{}"`
	Metadata  json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
