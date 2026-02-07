package domain

import (
	"encoding/json"
	"time"
)

// RentalReservationStatus represents the status of a reservation.
type RentalReservationStatus string

const (
	ReservationStatusPending            RentalReservationStatus = "ReservationPending"
	ReservationStatusConfirmed          RentalReservationStatus = "ReservationConfirmed"
	ReservationStatusCancelled          RentalReservationStatus = "ReservationCancelled"
	ReservationStatusPartiallyFulfilled RentalReservationStatus = "ReservationPartiallyFulfilled"
	ReservationStatusFulfilled          RentalReservationStatus = "ReservationFulfilled"
)

// RentalReservation represents an "intent to rent", aligning with schema.org/RentalReservation.
type RentalReservation struct {
	ID                int64                   `json:"id"`
	ReservationName   string                  `json:"reservationName,omitempty"` // schema.org/name
	ReservationStatus RentalReservationStatus `json:"reservationStatus"`         // schema.org/reservationStatus
	UnderNameID       *int64                  `json:"underNameId,omitempty"`     // Reference to Person
	BookingTime       time.Time               `json:"bookingTime"`               // When the reservation was made
	StartTime         time.Time               `json:"startTime"`                 // Expected start
	EndTime           time.Time               `json:"endTime"`                   // Expected end
	ProviderID        *int64                  `json:"providerId,omitempty"`      // Reference to Organization/Company
	Metadata          json.RawMessage         `json:"metadata,omitempty"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         time.Time               `json:"updatedAt"`

	Demands []Demand `json:"demands,omitempty"`
}

// Demand tracks the requirement for a specific type of asset, aligning with schema.org/Demand.
type Demand struct {
	ID               int64           `json:"id"`
	ReservationID    int64           `json:"reservationId,omitempty"`
	EventID          int64           `json:"eventId,omitempty"`
	ItemKind         string          `json:"itemKind"` // 'item_type', 'kit_template'
	ItemID           int64           `json:"itemId"`
	Quantity         int             `json:"requestedQuantity"`
	BusinessFunction string          `json:"businessFunction,omitempty"` // schema.org/businessFunction
	EligibleDuration string          `json:"eligibleDuration,omitempty"` // schema.org/eligibleDuration
	PlaceID          *int64          `json:"placeId,omitempty"`          // Where the demand is located
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

// CheckOutAction tracks the physical movement of assets out of the warehouse.
type CheckOutAction struct {
	ID            int64           `json:"id"`
	ReservationID int64           `json:"reservationId"`
	AssetID       int64           `json:"assetId"`
	AgentID       int64           `json:"agentId"`        // User/Person performing the action
	RecipientID   *int64          `json:"recipientId"`    // Person receiving the asset
	StartTime     time.Time       `json:"startTime"`      // When the checkout happened
	FromLocation  *int64          `json:"fromLocationId"` // PlaceID (Warehouse)
	ToLocation    *int64          `json:"toLocationId"`   // PlaceID (Event Location)
	Status        string          `json:"actionStatus"`   // schema.org/actionStatus
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

// ReturnAction tracks the physical movement of assets back to the warehouse.
type ReturnAction struct {
	ID            int64           `json:"id"`
	ReservationID int64           `json:"reservationId"`
	AssetID       int64           `json:"assetId"`
	AgentID       int64           `json:"agentId"`
	StartTime     time.Time       `json:"startTime"`
	FromLocation  *int64          `json:"fromLocationId"`
	ToLocation    *int64          `json:"toLocationId"` // PlaceID (Warehouse)
	Status        string          `json:"actionStatus"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

// FulfillmentLine represents the status of a single demand fulfillment.
type FulfillmentLine struct {
	DemandID          int64  `json:"demandId"`
	ItemKind          string `json:"itemKind"`
	ItemID            int64  `json:"itemId"`
	RequestedQuantity int    `json:"requestedQuantity"`
	FulfilledQuantity int    `json:"fulfilledQuantity"`
	ReturnedQuantity  int    `json:"returnedQuantity"`
	RemainingQuantity int    `json:"remainingQuantity"`
}

// RentalFulfillmentStatus summarizes the fulfillment state of a reservation.
type RentalFulfillmentStatus struct {
	ReservationID int64             `json:"reservationId"`
	Status        string            `json:"status"` // Overall status based on lines
	Lines         []FulfillmentLine `json:"lines"`
}
