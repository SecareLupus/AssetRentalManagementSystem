package domain

import (
	"encoding/json"
	"time"
)

type Company struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	LegalName   *string         `json:"legal_name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Contact struct {
	ID        int64           `json:"id"`
	CompanyID *int64          `json:"company_id,omitempty"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     *string         `json:"email,omitempty"`
	Phone     *string         `json:"phone,omitempty"`
	Role      *string         `json:"role,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Site struct {
	ID             int64           `json:"id"`
	CompanyID      int64           `json:"company_id"`
	Name           string          `json:"name"`
	AddressStreet  *string         `json:"address_street,omitempty"`
	AddressCity    *string         `json:"address_city,omitempty"`
	AddressState   *string         `json:"address_state,omitempty"`
	AddressZip     *string         `json:"address_zip,omitempty"`
	AddressCountry *string         `json:"address_country,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type Location struct {
	ID                 int64           `json:"id"`
	SiteID             int64           `json:"site_id"`
	ParentID           *int64          `json:"parent_id,omitempty"`
	Name               string          `json:"name"`
	LocationType       *string         `json:"location_type,omitempty"`
	PresumedAssetNeeds json.RawMessage `json:"presumed_asset_needs,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type EventStatus string

const (
	EventStatusAssumed   EventStatus = "assumed"
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID              int64           `json:"id"`
	CompanyID       int64           `json:"company_id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Status          EventStatus     `json:"status"`
	ParentEventID   *int64          `json:"parent_event_id,omitempty"`
	RecurrenceRule  *string         `json:"recurrence_rule,omitempty"`
	LastConfirmedAt *time.Time      `json:"last_confirmed_at,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type EventAssetNeed struct {
	ID         int64           `json:"id"`
	EventID    int64           `json:"event_id"`
	ItemTypeID int64           `json:"item_type_id"`
	Quantity   int             `json:"quantity"`
	IsAssumed  bool            `json:"is_assumed"`
	LocationID *int64          `json:"location_id,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
