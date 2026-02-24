package domain

import (
	"encoding/json"
	"time"
)

// ShowCompany acts as an overlay on top of the ingested domain.Company
type ShowCompany struct {
	ID        int64           `json:"id"`
	CompanyID int64           `json:"company_id"` // Refers to root domain.Company
	Metadata  json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Season belongs to a ShowCompany and spans a functional time period
type Season struct {
	ID            int64           `json:"id"`
	ShowCompanyID int64           `json:"show_company_id"`
	Name          string          `json:"name"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       time.Time       `json:"end_date"`
	Metadata      json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`

	Shows []Show `json:"shows,omitempty"`
}

// Ring represents a physical location/arena at a venue or ShowCompany
type Ring struct {
	ID            int64           `json:"id"`
	ShowCompanyID int64           `json:"show_company_id"` // The owning company/venue
	Name          string          `json:"name"`
	Description   *string         `json:"description,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Show is a specific event within a Season
type Show struct {
	ID         int64           `json:"id"`
	SeasonID   int64           `json:"season_id"`
	Name       string          `json:"name"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	LocationID *int64          `json:"location_id,omitempty"` // Optional link to a Place
	Metadata   json.RawMessage `json:"metadata,omitempty" swaggertype:"string" example:"{}"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`

	Rings []ShowRing `json:"rings,omitempty"`
}

// ShowRing maps a Ring to a specific Show and acts as the anchor for the loadout
type ShowRing struct {
	ID     int64 `json:"id"`
	ShowID int64 `json:"show_id"`
	RingID int64 `json:"ring_id"`

	Ring         *Ring             `json:"ring,omitempty"`
	LoadoutItems []RingLoadoutItem `json:"loadout_items,omitempty"`
}

// RingLoadoutItem defines a single requirement quantity for a specific Show's Ring
type RingLoadoutItem struct {
	ID         int64 `json:"id"`
	ShowRingID int64 `json:"show_ring_id"`
	ItemTypeID int64 `json:"item_type_id"`
	Quantity   int   `json:"quantity"`
}
