package domain

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	EventAssetCreated        EventType = "asset.created"
	EventAssetUpdated        EventType = "asset.updated"
	EventAssetTransitioned   EventType = "asset.status_changed"
	EventRentalSubmitted     EventType = "rental.submitted"
	EventRentalApproved      EventType = "rental.approved"
	EventItemTypeCreated     EventType = "item_type.created"
	EventInspectionSubmitted EventType = "inspection.submitted"
	EventInspectionSummary   EventType = "inspection.completed"
)

type OutboxStatus string

const (
	OutboxPending   OutboxStatus = "pending"
	OutboxProcessed OutboxStatus = "processed"
	OutboxFailed    OutboxStatus = "failed"
)

type OutboxEvent struct {
	ID           int64           `json:"id"`
	Type         EventType       `json:"event_type"`
	Payload      json.RawMessage `json:"payload"`
	Status       OutboxStatus    `json:"status"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	RetryCount   int             `json:"retry_count"`
	CreatedAt    time.Time       `json:"created_at"`
	ProcessedAt  *time.Time      `json:"processed_at,omitempty"`
}
