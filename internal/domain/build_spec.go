package domain

import (
	"encoding/json"
	"time"
)

type BuildSpec struct {
	ID             int64           `json:"id"`
	Version        string          `json:"version"`
	HardwareConfig json.RawMessage `json:"hardware_config,omitempty"`
	SoftwareConfig json.RawMessage `json:"software_config,omitempty"`
	FirmwareURL    *string         `json:"firmware_url,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
