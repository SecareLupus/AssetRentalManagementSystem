package worker

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/desmond/rental-management-system/internal/mqtt"
)

type OutboxWorker struct {
	repo       db.Repository
	httpClient *http.Client
	mqttClient *mqtt.Client
}

func NewOutboxWorker(repo db.Repository, mqttClient *mqtt.Client) *OutboxWorker {
	return &OutboxWorker{
		repo:       repo,
		mqttClient: mqttClient,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
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
	webhooks, err := w.repo.ListWebhooks(ctx)
	if err != nil {
		return fmt.Errorf("failed to list webhooks: %w", err)
	}

	for _, wh := range webhooks {
		// Check if webhook is interested in this event type
		interested := false
		for _, et := range wh.Events {
			if et == string(event.Type) || et == "*" {
				interested = true
				break
			}
		}

		if interested {
			if err := w.dispatchToWebhook(ctx, wh, event); err != nil {
				log.Printf("Webhook dispatch failed for URL %s: %v", wh.URL, err)
				// We log and continue to other webhooks.
				// In a more robust system, we might track delivery status per webhook.
			}
		}
	}

	// Mirror to MQTT
	if w.mqttClient != nil {
		topic := fmt.Sprintf("rms/events/%s", event.Type)
		payload, _ := json.Marshal(event)
		if err := w.mqttClient.Publish(topic, 1, false, payload); err != nil {
			log.Printf("MQTT Publish failed: %v", err)
		}
	}

	return nil
}

func (w *OutboxWorker) dispatchToWebhook(ctx context.Context, wh domain.WebhookConfig, event domain.OutboxEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", wh.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RMS-Event", string(event.Type))
	req.Header.Set("X-RMS-Delivery-ID", fmt.Sprintf("%d", event.ID))

	if wh.Secret != nil && *wh.Secret != "" {
		h := hmac.New(sha256.New, []byte(*wh.Secret))
		h.Write(body)
		signature := hex.EncodeToString(h.Sum(nil))
		req.Header.Set("X-RMS-Signature", signature)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
