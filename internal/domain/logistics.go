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

// DeliveryStatus represents the status of a shipment, aligning with schema.org/DeliveryStatusEvent.
type DeliveryStatus string

const (
	DeliveryPreparing DeliveryStatus = "DeliveryPreparing"
	DeliveryShipped   DeliveryStatus = "DeliveryShipped"
	DeliveryDelivered DeliveryStatus = "DeliveryDelivered"
	DeliveryReturned  DeliveryStatus = "DeliveryReturned"
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
	ID                  int64           `json:"id"`
	ReservationID       int64           `json:"reservationId"`
	AssetID             int64           `json:"assetId"`
	AgentID             int64           `json:"agentId"`             // User/Person performing the action
	RecipientID         *int64          `json:"recipientId"`         // Person receiving the asset
	ShipmentID          *int64          `json:"shipmentId"`          // Link to Shipment
	ScheduledDeliveryID *int64          `json:"scheduledDeliveryId"` // Link to ScheduledDelivery (if no shipment)
	StartTime           time.Time       `json:"startTime"`           // When the checkout happened
	FromLocation        *int64          `json:"fromLocationId"`      // PlaceID (Warehouse)
	ToLocation          *int64          `json:"toLocationId"`        // PlaceID (Event Location)
	Status              string          `json:"actionStatus"`        // schema.org/actionStatus
	Metadata            json.RawMessage `json:"metadata,omitempty"`
}

// ReturnAction tracks the physical movement of assets back to the warehouse.
type ReturnAction struct {
	ID            int64           `json:"id"`
	ReservationID int64           `json:"reservationId"`
	AssetID       int64           `json:"assetId"`
	AgentID       int64           `json:"agentId"`
	ShipmentID    *int64          `json:"shipmentId"` // Link to Shipment (e.g., return shipment)
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

// ScheduledDelivery represents a planned delivery for a specific Event/SeasonPlan (schema.org/DeliveryEvent).
type ScheduledDelivery struct {
	ID         int64     `json:"id"`
	EventID    int64     `json:"eventId"` // Replacing seasonPlanId for generic use
	TargetDate time.Time `json:"availableFrom"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ScheduledDeliveryItem represents the required items for a scheduled delivery (schema.org/Demand).
type ScheduledDeliveryItem struct {
	ID                  int64  `json:"id"`
	ScheduledDeliveryID int64  `json:"scheduledDeliveryId"`
	ItemKind            string `json:"itemKind"` // e.g., 'item_type'
	ItemID              int64  `json:"itemId"`
	Quantity            int    `json:"requestedQuantity"`
}

// Shipment represents a physical shipment, optionally linked to a ScheduledDelivery (schema.org/ParcelDelivery).
type Shipment struct {
	ID                  int64          `json:"id"`
	ScheduledDeliveryID *int64         `json:"scheduledDeliveryId,omitempty"`
	ProviderID          int64          `json:"providerId"` // Replacing showCompanyId
	ShipDate            time.Time      `json:"expectedArrivalFrom"`
	Carrier             string         `json:"carrier,omitempty"`
	TrackingNumber      string         `json:"trackingNumber,omitempty"`
	Status              DeliveryStatus `json:"deliveryStatus"` // schema.org/deliveryStatus (links to DeliveryStatusEvent)
	Notes               string         `json:"notes,omitempty"`
	Direction           string         `json:"direction"` // 'outbound' or 'inbound'
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
}
