package domain

type DashboardStats struct {
	TotalAssets       int            `json:"total_assets"`
	AssetsByStatus    map[string]int `json:"assets_by_status"`
	ActiveRentals     int            `json:"active_rentals"`
	PendingOutbox     int            `json:"pending_outbox_events"`
	RecentAlertsCount int            `json:"recent_alerts_count"`
}
