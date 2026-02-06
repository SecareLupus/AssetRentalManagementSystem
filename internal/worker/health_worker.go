package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
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
	liveTicker := time.NewTicker(time.Second * 10) // 10 second live pulses
	defer ticker.Stop()
	defer liveTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.CheckAllAssetsHealth(ctx)
		case <-liveTicker.C:
			w.PerformLivePolling(ctx)
		}
	}
}

func (w *HealthWorker) PerformLivePolling(ctx context.Context) {
	// In a real system, we'd only poll assets that have "Live Monitoring" enabled
	// or are currently Deployed/In-Use.
	assets, err := w.repo.ListAssets(ctx)
	if err != nil {
		return
	}

	for _, a := range assets {
		if a.Status != domain.AssetStatusDeployed || a.RemoteManagementID == nil || *a.RemoteManagementID == "" {
			continue
		}

		// Resolve manager from registry.
		// In a real system, the provider name would be in asset metadata or a 'provider' field.
		provider := "mock-provider"
		if len(a.Metadata) > 0 {
			var meta map[string]interface{}
			json.Unmarshal(a.Metadata, &meta)
			if p, ok := meta["remote_provider"].(string); ok {
				provider = p
			}
		}

		mgr, err := w.registry.Get(provider)
		if err != nil {
			continue
		}

		pulse, err := mgr.GetDevicePulse(ctx, *a.RemoteManagementID)
		if err != nil {
			tag := "unknown"
			if a.AssetTag != nil {
				tag = *a.AssetTag
			}
			log.Printf("HealthWorker: Failed to get pulse for asset %s via %s: %v", tag, provider, err)
			continue
		}

		tag := "unknown"
		if a.AssetTag != nil {
			tag = *a.AssetTag
		}

		topic := fmt.Sprintf("rms/assets/%s/pulse", tag)
		payload := fmt.Sprintf(`{"asset_tag": "%s", "pulse": %.2f, "timestamp": "%s"}`, tag, pulse, time.Now().Format(time.RFC3339))
		if w.mqttClient != nil {
			_ = w.mqttClient.Publish(topic, 0, false, payload)
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
