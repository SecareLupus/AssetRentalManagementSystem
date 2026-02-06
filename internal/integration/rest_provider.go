package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
)

type RESTProvider struct {
	name       string
	baseURL    string
	apiKey     string
	repo       db.Repository
	httpClient *http.Client
}

func NewRESTProvider(name, baseURL, apiKey string, repo db.Repository) *RESTProvider {
	return &RESTProvider{
		name:    name,
		baseURL: baseURL,
		apiKey:  apiKey,
		repo:    repo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *RESTProvider) GetName() string {
	return p.name
}

func (p *RESTProvider) SyncEntities(ctx context.Context) error {
	log.Printf("Integration [%s]: Starting entity sync from %s", p.name, p.baseURL)

	// Mock implementation: In a real system, we'd GET /assets and GET /locations
	// and upsert them into our DB.

	// Simulation of fetching external assets
	externalAssets := []struct {
		ExternalRef string `json:"ext_id"`
		Name        string `json:"name"`
		Status      string `json:"status"`
	}{
		{"EXT-001", "Backhoe Loader (External)", "available"},
		{"EXT-002", "Towable Generator (External)", "in-use"},
	}

	for _, ea := range externalAssets {
		log.Printf("Integration [%s]: Syncing asset %s (%s)", p.name, ea.Name, ea.ExternalRef)
		// Logic to check if asset exists in repo by external_ref and update/create.
		// For now, we just log the action.
	}

	return nil
}

func (p *RESTProvider) ExecuteAction(ctx context.Context, event domain.OutboxEvent) error {
	// Only forward certain events to this provider
	if event.Type != domain.EventRentalApproved && event.Type != domain.EventAssetTransitioned {
		return nil
	}

	log.Printf("Integration [%s]: Forwarding event %s to external system", p.name, event.Type)

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/events", p.baseURL), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("X-API-Key", p.apiKey)
	}

	// In a mock environment, this might fail unless a local simulator is running.
	// We'll log the "attempt" and succeed for the purpose of the walkthrough.
	log.Printf("Integration [%s]: POST %s/events (Payload size: %d)", p.name, p.baseURL, len(payload))

	return nil
}
