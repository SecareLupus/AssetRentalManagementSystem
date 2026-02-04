package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
)

type Repository interface {
	// ItemTypes
	CreateItemType(ctx context.Context, it *domain.ItemType) error
	GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error)
	ListItemTypes(ctx context.Context) ([]domain.ItemType, error)
	UpdateItemType(ctx context.Context, it *domain.ItemType) error
	DeleteItemType(ctx context.Context, id int64) error

	// Assets
	CreateAsset(ctx context.Context, a *domain.Asset) error
	GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error)
	ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error)
	UpdateAsset(ctx context.Context, a *domain.Asset) error
	UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus) error
	RecallAssetsByItemType(ctx context.Context, itemTypeID int64) error
	DeleteAsset(ctx context.Context, id int64) error

	// RentActions
	CreateRentAction(ctx context.Context, ra *domain.RentAction) error
	GetRentActionByID(ctx context.Context, id int64) (*domain.RentAction, error)
	UpdateRentAction(ctx context.Context, ra *domain.RentAction) error
	UpdateRentActionStatus(ctx context.Context, id int64, status domain.RentActionStatus, timestampField string, timestampValue time.Time) error

	// Inventory/Availability
	GetAvailableQuantity(ctx context.Context, itemTypeID int64, startTime, endTime time.Time) (int, error)

	// Maintenance
	AddMaintenanceLog(ctx context.Context, log *domain.MaintenanceLog) error
	ListMaintenanceLogs(ctx context.Context, assetID int64) ([]domain.MaintenanceLog, error)

	// Dynamic Inspections
	CreateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error
	GetInspectionTemplatesForItemType(ctx context.Context, itemTypeID int64) ([]domain.InspectionTemplate, error)
	SubmitInspection(ctx context.Context, is *domain.InspectionSubmission) error

	// Build Specs
	CreateBuildSpec(ctx context.Context, bs *domain.BuildSpec) error
	GetBuildSpecByID(ctx context.Context, id int64) (*domain.BuildSpec, error)
	ListBuildSpecs(ctx context.Context) ([]domain.BuildSpec, error)

	// Provisioning
	StartProvisioning(ctx context.Context, assetID int64, buildSpecID int64, performedBy string) (*domain.ProvisionAction, error)
	CompleteProvisioning(ctx context.Context, actionID int64, notes string) error

	// Users
	CreateUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, u *domain.User) error

	// Outbox
	AppendEvent(ctx context.Context, tx *sql.Tx, event *domain.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, id int64) error
	MarkEventFailed(ctx context.Context, id int64, errMessage string) error
}
