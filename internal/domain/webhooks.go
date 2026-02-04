package domain

import (
	"time"
)

type WebhookConfig struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	Secret    *string   `json:"secret,omitempty"`
	Events    []string  `json:"enabled_events"` // List of EventType strings
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
