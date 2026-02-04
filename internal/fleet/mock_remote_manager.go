package fleet

import (
	"context"
	"fmt"

	"github.com/desmond/rental-management-system/internal/domain"
)

type MockRemoteManager struct {
	DeviceHealth map[string]domain.RemoteHealthStatus
}

func NewMockRemoteManager() *MockRemoteManager {
	return &MockRemoteManager{
		DeviceHealth: make(map[string]domain.RemoteHealthStatus),
	}
}

func (m *MockRemoteManager) GetDeviceInfo(ctx context.Context, remoteID string) (*domain.DeviceInfo, error) {
	health, ok := m.DeviceHealth[remoteID]
	if !ok {
		health = domain.HealthOnline // Default for mock
	}

	return &domain.DeviceInfo{
		RemoteID:     remoteID,
		HealthStatus: health,
		IPAddress:    "192.168.1.100",
		AgentVersion: "mock-v1.0",
	}, nil
}

func (m *MockRemoteManager) ApplyPowerAction(ctx context.Context, remoteID string, action domain.RemotePowerAction) error {
	fmt.Printf("Mock: Applying %s to device %s\n", action, remoteID)
	return nil
}

func (m *MockRemoteManager) GetName() string {
	return "mock-provider"
}
func (m *MockRemoteManager) GetManager() domain.RemoteManager {
	return m
}
