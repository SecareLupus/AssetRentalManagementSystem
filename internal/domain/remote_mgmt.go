package domain

import (
	"context"
)

type RemotePowerAction string

const (
	PowerOn     RemotePowerAction = "on"
	PowerOff    RemotePowerAction = "off"
	PowerReboot RemotePowerAction = "reboot"
)

type RemoteHealthStatus string

const (
	HealthOnline  RemoteHealthStatus = "online"
	HealthOffline RemoteHealthStatus = "offline"
	HealthUnknown RemoteHealthStatus = "unknown"
)

type DeviceInfo struct {
	RemoteID     string             `json:"remote_id"`
	HealthStatus RemoteHealthStatus `json:"health_status"`
	Uptime       int64              `json:"uptime,omitempty"`
	IPAddress    string             `json:"ip_address,omitempty"`
	AgentVersion string             `json:"agent_version,omitempty"`
}

type RemoteManager interface {
	GetDeviceInfo(ctx context.Context, remoteID string) (*DeviceInfo, error)
	ApplyPowerAction(ctx context.Context, remoteID string, action RemotePowerAction) error
}

type RemoteManagerProvider interface {
	GetName() string
	GetManager() RemoteManager
}
