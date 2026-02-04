package domain

import (
	"time"
)

type ProvisionStatus string

const (
	ProvisionStarted   ProvisionStatus = "started"
	ProvisionCompleted ProvisionStatus = "completed"
	ProvisionFailed    ProvisionStatus = "failed"
)

type ProvisionAction struct {
	ID          int64           `json:"id"`
	AssetID     int64           `json:"asset_id"`
	BuildSpecID *int64          `json:"build_spec_id,omitempty"`
	Status      ProvisionStatus `json:"status"`
	PerformedBy string          `json:"performed_by"`
	Notes       *string         `json:"notes,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}
