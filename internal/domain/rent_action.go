package domain

import (
	"encoding/json"
	"fmt"
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
	ID              int64            `json:"id"`
	RequesterRef    string           `json:"requester_ref"`
	CreatedByRef    string           `json:"created_by_ref"`               // Legacy ref or username
	CreatedByUserID *int64           `json:"created_by_user_id,omitempty"` // Audit trail
	UpdatedByUserID *int64           `json:"updated_by_user_id,omitempty"`
	ApprovedByRef   *string          `json:"approved_by_ref,omitempty"`
	Status          RentActionStatus `json:"status"`
	Priority        string           `json:"priority"`
	StartTime       time.Time        `json:"start_time"` // Maps to schema.org/startTime
	EndTime         time.Time        `json:"end_time"`   // Maps to schema.org/endTime
	IsASAP          bool             `json:"is_asap"`
	Description     *string          `json:"description,omitempty"` // Maps to schema.org/description
	ExternalSource  *string          `json:"external_source,omitempty"`
	ExternalRef     *string          `json:"external_ref,omitempty"`
	SchemaOrg       json.RawMessage  `json:"schema_org,omitempty" swaggertype:"string" example:"{}"`
	Metadata        json.RawMessage  `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	ApprovedAt      *time.Time       `json:"approved_at,omitempty"`
	RejectedAt      *time.Time       `json:"rejected_at,omitempty"`
	CancelledAt     *time.Time       `json:"cancelled_at,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`

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
	Metadata          json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
}

func (ra *RentAction) Validate() error {
	if ra.StartTime.After(ra.EndTime) || ra.StartTime.Equal(ra.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}
	if ra.RequesterRef == "" {
		return fmt.Errorf("requester is required")
	}
	return nil
}

func (ra *RentAction) Submit() error {
	if ra.Status != RentActionStatusDraft {
		return fmt.Errorf("can only submit from draft status, current: %s", ra.Status)
	}
	ra.Status = RentActionStatusPending
	return nil
}

func (ra *RentAction) Approve() error {
	if ra.Status != RentActionStatusPending {
		return fmt.Errorf("can only approve from pending status, current: %s", ra.Status)
	}
	now := time.Now()
	ra.Status = RentActionStatusApproved
	ra.ApprovedAt = &now
	return nil
}

func (ra *RentAction) Reject() error {
	if ra.Status != RentActionStatusPending {
		return fmt.Errorf("can only reject from pending status, current: %s", ra.Status)
	}
	now := time.Now()
	ra.Status = RentActionStatusRejected
	ra.RejectedAt = &now
	return nil
}

func (ra *RentAction) Cancel() error {
	if ra.Status == RentActionStatusFulfilled || ra.Status == RentActionStatusCancelled {
		return fmt.Errorf("cannot cancel reservation in status: %s", ra.Status)
	}
	now := time.Now()
	ra.Status = RentActionStatusCancelled
	ra.CancelledAt = &now
	return nil
}
