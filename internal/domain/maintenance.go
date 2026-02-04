package domain

import (
	"encoding/json"
	"time"
)

type MaintenanceActionType string

const (
	ActionInspect   MaintenanceActionType = "inspect"
	ActionRepair    MaintenanceActionType = "repair"
	ActionUpgrade   MaintenanceActionType = "upgrade"
	ActionRefurbish MaintenanceActionType = "refurbish"
)

type MaintenanceLog struct {
	ID          int64                 `json:"id"`
	AssetID     int64                 `json:"asset_id"`
	ActionType  MaintenanceActionType `json:"action_type"`
	Notes       *string               `json:"notes,omitempty"`
	PerformedBy string                `json:"performed_by"`
	TestBits    json.RawMessage       `json:"test_bits,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
}
