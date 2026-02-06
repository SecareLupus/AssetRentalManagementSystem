package integration

import (
	"context"

	"github.com/desmond/rental-management-system/internal/domain"
)

// ExternalProvider defines the interface for third-party system integrations.
type ExternalProvider interface {
	GetName() string
	// SyncEntities pulls data from the external system (Assets, Locations, etc.)
	SyncEntities(ctx context.Context) error
	// ExecuteAction pushes an event-driven action to the external system.
	ExecuteAction(ctx context.Context, event domain.OutboxEvent) error
}

// IntegrationService manages a collection of external providers.
type IntegrationService struct {
	providers []ExternalProvider
}

func NewIntegrationService() *IntegrationService {
	return &IntegrationService{
		providers: make([]ExternalProvider, 0),
	}
}

func (s *IntegrationService) RegisterProvider(p ExternalProvider) {
	s.providers = append(s.providers, p)
}

func (s *IntegrationService) SyncAll(ctx context.Context) {
	for _, p := range s.providers {
		if err := p.SyncEntities(ctx); err != nil {
			// In a real system, we'd log this or track it in a status table
			_ = err
		}
	}
}

func (s *IntegrationService) HandleEvent(ctx context.Context, event domain.OutboxEvent) error {
	for _, p := range s.providers {
		if err := p.ExecuteAction(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
