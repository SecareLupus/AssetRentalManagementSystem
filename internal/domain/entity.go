package domain

import (
	"encoding/json"
	"time"
)

// Schema.org basic types

type PostalAddress struct {
	StreetAddress   *string `json:"street_address,omitempty"`
	AddressLocality *string `json:"address_locality,omitempty"` // City
	AddressRegion   *string `json:"address_region,omitempty"`   // State/Region
	PostalCode      *string `json:"postal_code,omitempty"`
	AddressCountry  *string `json:"address_country,omitempty"`
}

type ContactPoint struct {
	Email string  `json:"email,omitempty"`
	Phone string  `json:"phone,omitempty"`
	Type  *string `json:"type,omitempty"` // e.g., "customer service", "technical support"
}

type Company struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	LegalName   *string         `json:"legal_name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
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
	PlaceID    *int64          `json:"place_id,omitempty"` // Replaces LocationID
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// Phase 24 Additions

type Place struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	Description        *string         `json:"description,omitempty"`
	ContainedInPlaceID *int64          `json:"contained_in_place_id,omitempty"`
	OwnerID            *int64          `json:"owner_id,omitempty"` // Company/Org owning the root place
	Category           *string         `json:"category,omitempty"` // "site", "room", "zone", etc.
	Address            *PostalAddress  `json:"address,omitempty"`
	IsInternal         bool            `json:"is_internal"`
	PresumedDemands    json.RawMessage `json:"presumed_demands,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type Person struct {
	ID            int64           `json:"id"`
	GivenName     string          `json:"given_name"`
	FamilyName    string          `json:"family_name"`
	ContactPoints []ContactPoint  `json:"contact_points,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type OrganizationRole struct {
	ID             int64           `json:"id"`
	PersonID       int64           `json:"person_id"`
	OrganizationID int64           `json:"organization_id"`
	RoleName       string          `json:"role_name"` // "Contact", "Manager", etc.
	StartDate      *time.Time      `json:"start_date,omitempty"`
	EndDate        *time.Time      `json:"end_date,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}
