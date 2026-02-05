package domain

import "time"

type AvailabilityPoint struct {
	Date      time.Time `json:"date"`
	Available int       `json:"available"`
	Total     int       `json:"total"`
}

type ShortageAlert struct {
	ItemTypeID    int64     `json:"item_type_id"`
	ItemTypeName  string    `json:"item_type_name"`
	Date          time.Time `json:"date"`
	ShortageCount int       `json:"shortage_count"`
	TotalNeeded   int       `json:"total_needed"`
	TotalOwned    int       `json:"total_owned"`
}

type MaintenanceForecast struct {
	AssetID      int64     `json:"asset_id"`
	AssetTag     string    `json:"asset_tag"`
	NextService  time.Time `json:"next_service_date"`
	Reason       string    `json:"reason"`
	UrgencyScore float64   `json:"urgency_score"` // 0-1
}
