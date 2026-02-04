package worker

import (
	"context"
	"log"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
)

type OutboxWorker struct {
	repo db.Repository
}

func NewOutboxWorker(repo db.Repository) *OutboxWorker {
	return &OutboxWorker{repo: repo}
}

func (w *OutboxWorker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.ProcessEvents(ctx)
		}
	}
}

func (w *OutboxWorker) ProcessEvents(ctx context.Context) {
	events, err := w.repo.GetPendingEvents(ctx, 10)
	if err != nil {
		log.Printf("Failed to fetch pending events: %v", err)
		return
	}

	for _, event := range events {
		if err := w.DeliverEvent(ctx, event); err != nil {
			log.Printf("Failed to deliver event %d: %v", event.ID, err)
			w.repo.MarkEventFailed(ctx, event.ID, err.Error())
		} else {
			w.repo.MarkEventProcessed(ctx, event.ID)
		}
	}
}

func (w *OutboxWorker) DeliverEvent(ctx context.Context, event domain.OutboxEvent) error {
	// For now, we simulate delivery by logging
	// In a real system, this would iterate through configured webhooks or external sync adapters
	log.Printf("Delivering event: %s (ID: %d), Payload: %s", event.Type, event.ID, string(event.Payload))

	// Real implementation example:
	// return w.webhookSvc.Dispatch(ctx, event)

	return nil
}
