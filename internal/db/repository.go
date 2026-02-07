package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
)

type Repository interface {
	// ItemTypes
	CreateItemType(ctx context.Context, it *domain.ItemType) error
	GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error)
	ListItemTypes(ctx context.Context, includeInactive bool) ([]domain.ItemType, error)
	UpdateItemType(ctx context.Context, it *domain.ItemType) error
	DeleteItemType(ctx context.Context, id int64) error

	// Assets
	CreateAsset(ctx context.Context, a *domain.Asset) error
	GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error)
	ListAssets(ctx context.Context) ([]domain.Asset, error)
	ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error)
	UpdateAsset(ctx context.Context, a *domain.Asset) error
	UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus, placeID *int64, location *string, metadata json.RawMessage) error
	RecallAssetsByItemType(ctx context.Context, itemTypeID int64) error
	BulkRecallAssets(ctx context.Context, ids []int64) error
	DeleteAsset(ctx context.Context, id int64) error

	GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error)

	// Logistics
	CreateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error
	GetRentalReservationByID(ctx context.Context, id int64) (*domain.RentalReservation, error)
	ListRentalReservations(ctx context.Context) ([]domain.RentalReservation, error)
	UpdateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error
	UpdateRentalReservationStatus(ctx context.Context, id int64, status domain.RentalReservationStatus) error

	CreateDemand(ctx context.Context, d *domain.Demand) error
	ListDemandsByReservation(ctx context.Context, reservationID int64) ([]domain.Demand, error)
	ListDemandsByEvent(ctx context.Context, eventID int64) ([]domain.Demand, error)
	UpdateDemand(ctx context.Context, d *domain.Demand) error
	DeleteDemand(ctx context.Context, id int64) error

	CreateCheckOutAction(ctx context.Context, co *domain.CheckOutAction) error
	CreateReturnAction(ctx context.Context, ra *domain.ReturnAction) error
	ListCheckOutActions(ctx context.Context, reservationID int64) ([]domain.CheckOutAction, error)
	ListReturnActions(ctx context.Context, reservationID int64) ([]domain.ReturnAction, error)
	GetRentalFulfillmentStatus(ctx context.Context, reservationID int64) (*domain.RentalFulfillmentStatus, error)
	BatchCheckOut(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64, fromLocationID, toLocationID *int64) error
	BatchReturn(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64, toLocationID *int64) error

	// Inventory/Availability
	GetAvailableQuantity(ctx context.Context, itemTypeID int64, startTime, endTime time.Time) (int, error)

	// Maintenance
	AddMaintenanceLog(ctx context.Context, log *domain.MaintenanceLog) error
	ListMaintenanceLogs(ctx context.Context, assetID int64) ([]domain.MaintenanceLog, error)

	// Maintenance & Inspections
	CreateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error
	UpdateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error
	DeleteInspectionTemplate(ctx context.Context, id int64) error
	ListInspectionTemplates(ctx context.Context) ([]domain.InspectionTemplate, error)
	GetInspectionTemplate(ctx context.Context, id int64) (*domain.InspectionTemplate, error)
	GetInspectionTemplatesForItemType(ctx context.Context, itemTypeID int64) ([]domain.InspectionTemplate, error)
	SetItemTypeInspections(ctx context.Context, itemTypeID int64, templateIDs []int64) error
	CreateInspectionSubmission(ctx context.Context, is *domain.InspectionSubmission) error

	// Entity Management
	CreateCompany(ctx context.Context, c *domain.Company) error
	GetCompany(ctx context.Context, id int64) (*domain.Company, error)
	ListCompanies(ctx context.Context) ([]domain.Company, error)
	UpdateCompany(ctx context.Context, c *domain.Company) error
	DeleteCompany(ctx context.Context, id int64) error

	// Unified Person & Role Management
	CreatePerson(ctx context.Context, p *domain.Person) error
	GetPerson(ctx context.Context, id int64) (*domain.Person, error)
	ListPeople(ctx context.Context) ([]domain.Person, error)
	UpdatePerson(ctx context.Context, p *domain.Person) error
	DeletePerson(ctx context.Context, id int64) error

	CreateOrganizationRole(ctx context.Context, or *domain.OrganizationRole) error
	ListOrganizationRoles(ctx context.Context, orgID *int64, personID *int64) ([]domain.OrganizationRole, error)
	DeleteOrganizationRole(ctx context.Context, id int64) error

	// Unified Place Management
	CreatePlace(ctx context.Context, p *domain.Place) error
	GetPlace(ctx context.Context, id int64) (*domain.Place, error)
	ListPlaces(ctx context.Context, ownerID *int64, parentID *int64) ([]domain.Place, error)
	UpdatePlace(ctx context.Context, p *domain.Place) error
	DeletePlace(ctx context.Context, id int64) error

	CreateEvent(ctx context.Context, e *domain.Event) error
	GetEvent(ctx context.Context, id int64) (*domain.Event, error)
	ListEvents(ctx context.Context, companyID *int64) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, e *domain.Event) error
	DeleteEvent(ctx context.Context, id int64) error

	// Demands (event context)
	// These are also handled via CreateDemand/ListDemandsByEvent

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

	// Intelligence
	GetAvailabilityTimeline(ctx context.Context, itemTypeID int64, start, end time.Time) ([]domain.AvailabilityPoint, error)
	GetShortageAlerts(ctx context.Context) ([]domain.ShortageAlert, error)
	GetMaintenanceForecast(ctx context.Context) ([]domain.MaintenanceForecast, error)

	// Outbox / Webhooks
	AppendEvent(ctx context.Context, tx *sql.Tx, event *domain.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, id int64) error
	MarkEventFailed(ctx context.Context, id int64, errMessage string) error
	ListWebhooks(ctx context.Context) ([]domain.WebhookConfig, error)
}
