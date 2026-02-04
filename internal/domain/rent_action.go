package domain

import (
	"encoding/json"
	"time"
)

type RentActionStatus string

const (
	RentActionStatusDraft     RentActionStatus = "draft"
	RentActionStatusPending   RentActionStatus = "pending"
	RentActionStatusApproved  RentActionStatus = "approved"
	RentActionStatusRejected  RentActionStatus = "rejected"
	RentActionStatusCancelled RentActionStatus = "cancelled"
	RentActionStatusFulfilled RentActionStatus = "fulfilled"
)

type RentAction struct {
	ID             int64            `json:"id"`
	RequesterRef   string           `json:"requester_ref"`
	CreatedByRef   string           `json:"created_by_ref"`
	ApprovedByRef  *string          `json:"approved_by_ref,omitempty"`
	Status         RentActionStatus `json:"status"`
	Priority       string           `json:"priority"`
	StartTime      time.Time        `json:"start_time"` // Maps to schema.org/startTime
	EndTime        time.Time        `json:"end_time"`   // Maps to schema.org/endTime
	IsASAP         bool             `json:"is_asap"`
	Description    *string          `json:"description,omitempty"` // Maps to schema.org/description
	ExternalSource *string          `json:"external_source,omitempty"`
	ExternalRef    *string          `json:"external_ref,omitempty"`
	SchemaOrg      json.RawMessage  `json:"schema_org,omitempty"`
	Metadata       json.RawMessage  `json:"metadata,omitempty"`
	ApprovedAt     *time.Time       `json:"approved_at,omitempty"`
	RejectedAt     *time.Time       `json:"rejected_at,omitempty"`
	CancelledAt    *time.Time       `json:"cancelled_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`

	Items []RentActionItem `json:"items,omitempty"`
}

type RentActionItem struct {
	ID                int64           `json:"id"`
	RentActionID      int64           `json:"rent_action_id"`
	ItemKind          string          `json:"item_kind"` // 'item_type','kit_template'
	ItemID            int64           `json:"item_id"`
	RequestedQuantity int             `json:"requested_quantity"`
	AllocatedQuantity int             `json:"allocated_quantity"`
	Notes             *string         `json:"notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
}
