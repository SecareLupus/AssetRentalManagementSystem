package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/fleet"
	"github.com/desmond/rental-management-system/internal/mqtt"
)

type HealthWorker struct {
	repo       db.Repository
	mqttClient *mqtt.Client
	registry   *fleet.RemoteRegistry
}

func NewHealthWorker(repo db.Repository, mqttClient *mqtt.Client, registry *fleet.RemoteRegistry) *HealthWorker {
	return &HealthWorker{
		repo:       repo,
		mqttClient: mqttClient,
		registry:   registry,
	}
}

func (w *HealthWorker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.CheckAllAssetsHealth(ctx)
		}
	}
}

func (w *HealthWorker) CheckAllAssetsHealth(ctx context.Context) {
	assets, err := w.repo.ListAssets(ctx)
	if err != nil {
		log.Printf("HealthWorker: Failed to list assets: %v", err)
		return
	}

	for _, a := range assets {
		if a.RemoteManagementID == nil || *a.RemoteManagementID == "" {
			continue
		}

		// Assume provider is stored in metadata or default for now
		// In a real system, asset would have a provider field.
		// For Phase 13, we stick to the registration logic in fleet.

		// For now, we try to get a manager. If we don't have enough metadata, we skip.
		// Mock manager is usually registered as "mock-provider"
		mgr, err := w.registry.Get("mock-provider")
		if err != nil {
			continue
		}

		info, err := mgr.GetDeviceInfo(ctx, *a.RemoteManagementID)
		if err != nil {
			tag := "unknown"
			if a.AssetTag != nil {
				tag = *a.AssetTag
			}
			log.Printf("HealthWorker: Failed to get health for asset %s: %v", tag, err)
			continue
		}

		// Publish to MQTT
		tag := "unknown"
		if a.AssetTag != nil {
			tag = *a.AssetTag
		}
		topic := fmt.Sprintf("rms/assets/%s/health", tag)
		payload, _ := json.Marshal(info)
		if w.mqttClient != nil {
			if err := w.mqttClient.Publish(topic, 1, true, payload); err != nil {
				log.Printf("HealthWorker: MQTT publish failed: %v", err)
			}
		}
	}
}
