package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// Implement only what's used by ingestItem and upsertAsset
func (m *MockRepository) ListItemTypes(ctx context.Context, inc bool) ([]domain.ItemType, error) {
	args := m.Called(ctx, inc)
	return args.Get(0).([]domain.ItemType), args.Error(1)
}

func (m *MockRepository) UpsertItemType(ctx context.Context, it *domain.ItemType) error {
	args := m.Called(ctx, it)
	if it != nil && args.Get(0) != nil {
		it.ID = args.Get(0).(*domain.ItemType).ID
	}
	return args.Error(1)
}

func (m *MockRepository) UpsertAsset(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

// Dummy implementations for the rest of Repository interface
func (m *MockRepository) CreateItemType(ctx context.Context, it *domain.ItemType) error { return nil }
func (m *MockRepository) GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error) {
	return nil, nil
}
func (m *MockRepository) UpdateItemType(ctx context.Context, it *domain.ItemType) error { return nil }
func (m *MockRepository) DeleteItemType(ctx context.Context, id int64) error            { return nil }
func (m *MockRepository) CreateAsset(ctx context.Context, a *domain.Asset) error        { return nil }
func (m *MockRepository) GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error) {
	return nil, nil
}
func (m *MockRepository) ListAssets(ctx context.Context) ([]domain.Asset, error) { return nil, nil }
func (m *MockRepository) ListAssetsByItemType(ctx context.Context, id int64) ([]domain.Asset, error) {
	return nil, nil
}
func (m *MockRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error { return nil }
func (m *MockRepository) UpdateAssetStatus(ctx context.Context, id int64, s domain.AssetStatus, p *int64, l *string, mdt json.RawMessage) error {
	return nil
}
func (m *MockRepository) RecallAssetsByItemType(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) BulkRecallAssets(ctx context.Context, ids []int64) error    { return nil }
func (m *MockRepository) DeleteAsset(ctx context.Context, id int64) error            { return nil }
func (m *MockRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	return nil, nil
}
func (m *MockRepository) CreateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	return nil
}
func (m *MockRepository) GetRentalReservationByID(ctx context.Context, id int64) (*domain.RentalReservation, error) {
	return nil, nil
}
func (m *MockRepository) ListRentalReservations(ctx context.Context) ([]domain.RentalReservation, error) {
	return nil, nil
}
func (m *MockRepository) UpdateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	return nil
}
func (m *MockRepository) UpdateRentalReservationStatus(ctx context.Context, id int64, s domain.RentalReservationStatus) error {
	return nil
}
func (m *MockRepository) CreateDemand(ctx context.Context, d *domain.Demand) error { return nil }
func (m *MockRepository) ListDemandsByReservation(ctx context.Context, id int64) ([]domain.Demand, error) {
	return nil, nil
}
func (m *MockRepository) ListDemandsByEvent(ctx context.Context, id int64) ([]domain.Demand, error) {
	return nil, nil
}
func (m *MockRepository) UpdateDemand(ctx context.Context, d *domain.Demand) error { return nil }
func (m *MockRepository) DeleteDemand(ctx context.Context, id int64) error         { return nil }
func (m *MockRepository) CreateCheckOutAction(ctx context.Context, co *domain.CheckOutAction) error {
	return nil
}
func (m *MockRepository) CreateReturnAction(ctx context.Context, ra *domain.ReturnAction) error {
	return nil
}
func (m *MockRepository) ListCheckOutActions(ctx context.Context, id int64) ([]domain.CheckOutAction, error) {
	return nil, nil
}
func (m *MockRepository) ListReturnActions(ctx context.Context, id int64) ([]domain.ReturnAction, error) {
	return nil, nil
}
func (m *MockRepository) GetRentalFulfillmentStatus(ctx context.Context, id int64) (*domain.RentalFulfillmentStatus, error) {
	return nil, nil
}
func (m *MockRepository) BatchCheckOut(ctx context.Context, id int64, ids []int64, aid int64, f, t *int64) error {
	return nil
}
func (m *MockRepository) BatchReturn(ctx context.Context, id int64, ids []int64, aid int64, t *int64) error {
	return nil
}
func (m *MockRepository) GetAvailableQuantity(ctx context.Context, id int64, s, e time.Time) (int, error) {
	return 0, nil
}
func (m *MockRepository) AddMaintenanceLog(ctx context.Context, l *domain.MaintenanceLog) error {
	return nil
}
func (m *MockRepository) ListMaintenanceLogs(ctx context.Context, id int64) ([]domain.MaintenanceLog, error) {
	return nil, nil
}
func (m *MockRepository) CreateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	return nil
}
func (m *MockRepository) UpdateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	return nil
}
func (m *MockRepository) DeleteInspectionTemplate(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) ListInspectionTemplates(ctx context.Context) ([]domain.InspectionTemplate, error) {
	return nil, nil
}
func (m *MockRepository) GetInspectionTemplate(ctx context.Context, id int64) (*domain.InspectionTemplate, error) {
	return nil, nil
}
func (m *MockRepository) GetInspectionTemplatesForItemType(ctx context.Context, id int64) ([]domain.InspectionTemplate, error) {
	return nil, nil
}
func (m *MockRepository) SetItemTypeInspections(ctx context.Context, id int64, ids []int64) error {
	return nil
}
func (m *MockRepository) CreateInspectionSubmission(ctx context.Context, is *domain.InspectionSubmission) error {
	return nil
}
func (m *MockRepository) CreateCompany(ctx context.Context, c *domain.Company) error { return nil }
func (m *MockRepository) GetCompany(ctx context.Context, id int64) (*domain.Company, error) {
	return nil, nil
}
func (m *MockRepository) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	return nil, nil
}
func (m *MockRepository) UpdateCompany(ctx context.Context, c *domain.Company) error { return nil }
func (m *MockRepository) DeleteCompany(ctx context.Context, id int64) error          { return nil }
func (m *MockRepository) CreatePerson(ctx context.Context, p *domain.Person) error   { return nil }
func (m *MockRepository) GetPerson(ctx context.Context, id int64) (*domain.Person, error) {
	return nil, nil
}
func (m *MockRepository) ListPeople(ctx context.Context) ([]domain.Person, error)  { return nil, nil }
func (m *MockRepository) UpdatePerson(ctx context.Context, p *domain.Person) error { return nil }
func (m *MockRepository) DeletePerson(ctx context.Context, id int64) error         { return nil }
func (m *MockRepository) CreateOrganizationRole(ctx context.Context, or *domain.OrganizationRole) error {
	return nil
}
func (m *MockRepository) ListOrganizationRoles(ctx context.Context, oid, pid *int64) ([]domain.OrganizationRole, error) {
	return nil, nil
}
func (m *MockRepository) DeleteOrganizationRole(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) CreatePlace(ctx context.Context, p *domain.Place) error     { return nil }
func (m *MockRepository) GetPlace(ctx context.Context, id int64) (*domain.Place, error) {
	return nil, nil
}
func (m *MockRepository) ListPlaces(ctx context.Context, oid, pid *int64) ([]domain.Place, error) {
	return nil, nil
}
func (m *MockRepository) UpdatePlace(ctx context.Context, p *domain.Place) error { return nil }
func (m *MockRepository) DeletePlace(ctx context.Context, id int64) error        { return nil }
func (m *MockRepository) CreateEvent(ctx context.Context, e *domain.Event) error { return nil }
func (m *MockRepository) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) ListEvents(ctx context.Context, cid *int64) ([]domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) UpdateEvent(ctx context.Context, e *domain.Event) error          { return nil }
func (m *MockRepository) DeleteEvent(ctx context.Context, id int64) error                 { return nil }
func (m *MockRepository) CreateBuildSpec(ctx context.Context, bs *domain.BuildSpec) error { return nil }
func (m *MockRepository) GetBuildSpecByID(ctx context.Context, id int64) (*domain.BuildSpec, error) {
	return nil, nil
}
func (m *MockRepository) ListBuildSpecs(ctx context.Context) ([]domain.BuildSpec, error) {
	return nil, nil
}
func (m *MockRepository) StartProvisioning(ctx context.Context, aid, bid int64, pb string) (*domain.ProvisionAction, error) {
	return nil, nil
}
func (m *MockRepository) CompleteProvisioning(ctx context.Context, aid int64, n string) error {
	return nil
}
func (m *MockRepository) CreateUser(ctx context.Context, u *domain.User) error { return nil }
func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return nil, nil
}
func (m *MockRepository) GetUserByUsername(ctx context.Context, u string) (*domain.User, error) {
	return nil, nil
}
func (m *MockRepository) UpdateUser(ctx context.Context, u *domain.User) error { return nil }
func (m *MockRepository) ListUsers(ctx context.Context) ([]domain.User, error) { return nil, nil }
func (m *MockRepository) DeleteUser(ctx context.Context, id int64) error       { return nil }
func (m *MockRepository) GetSettings(ctx context.Context) (map[string]json.RawMessage, error) {
	return nil, nil
}
func (m *MockRepository) UpdateSetting(ctx context.Context, k string, v json.RawMessage) error {
	return nil
}
func (m *MockRepository) GetAvailabilityTimeline(ctx context.Context, id int64, s, e time.Time) ([]domain.AvailabilityPoint, error) {
	return nil, nil
}
func (m *MockRepository) GetShortageAlerts(ctx context.Context) ([]domain.ShortageAlert, error) {
	return nil, nil
}
func (m *MockRepository) GetMaintenanceForecast(ctx context.Context) ([]domain.MaintenanceForecast, error) {
	return nil, nil
}
func (m *MockRepository) CreateIngestSource(ctx context.Context, s *domain.IngestSource) error {
	return nil
}
func (m *MockRepository) UpdateIngestSource(ctx context.Context, s *domain.IngestSource) error {
	return nil
}
func (m *MockRepository) ListIngestSources(ctx context.Context) ([]domain.IngestSource, error) {
	return nil, nil
}
func (m *MockRepository) GetIngestSource(ctx context.Context, id int64) (*domain.IngestSource, error) {
	return nil, nil
}
func (m *MockRepository) DeleteIngestSource(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) GetPendingIngestSources(ctx context.Context) ([]domain.IngestSource, error) {
	return nil, nil
}
func (m *MockRepository) CreateIngestEndpoint(ctx context.Context, ep *domain.IngestEndpoint) error {
	return nil
}
func (m *MockRepository) GetIngestEndpoint(ctx context.Context, id int64) (*domain.IngestEndpoint, error) {
	return nil, nil
}
func (m *MockRepository) UpdateIngestEndpoint(ctx context.Context, ep *domain.IngestEndpoint) error {
	return nil
}
func (m *MockRepository) DeleteIngestEndpoint(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) ListIngestEndpoints(ctx context.Context, sid int64) ([]domain.IngestEndpoint, error) {
	return nil, nil
}
func (m *MockRepository) SetEndpointMappings(ctx context.Context, eid int64, ms []domain.IngestMapping) error {
	return nil
}
func (m *MockRepository) UpsertCompany(ctx context.Context, c *domain.Company) error { return nil }
func (m *MockRepository) UpsertPerson(ctx context.Context, p *domain.Person) error   { return nil }
func (m *MockRepository) UpsertPlace(ctx context.Context, p *domain.Place) error     { return nil }
func (m *MockRepository) AppendEvent(ctx context.Context, tx *sql.Tx, e *domain.OutboxEvent) error {
	return nil
}
func (m *MockRepository) GetPendingEvents(ctx context.Context, l int) ([]domain.OutboxEvent, error) {
	return nil, nil
}
func (m *MockRepository) MarkEventProcessed(ctx context.Context, id int64) error         { return nil }
func (m *MockRepository) MarkEventFailed(ctx context.Context, id int64, em string) error { return nil }
func (m *MockRepository) ListWebhooks(ctx context.Context) ([]domain.WebhookConfig, error) {
	return nil, nil
}

func TestIngestWorker_ItemTypeInference(t *testing.T) {
	repo := new(MockRepository)
	worker := NewIngestWorker(repo)
	ctx := context.Background()

	// 1. Initial Repo State: No item types
	repo.On("ListItemTypes", ctx, true).Return([]domain.ItemType{}, nil)

	// 2. Mappings
	mappings := []domain.IngestMapping{
		{JSONPath: "$.device_type", TargetModel: domain.IngestTargetItemType, TargetField: "name", IsIdentity: true},
		{JSONPath: "$.device_type", TargetModel: domain.IngestTargetAsset, TargetField: "item_type_name"},
		{JSONPath: "$.serial", TargetModel: domain.IngestTargetAsset, TargetField: "serial_number", IsIdentity: true},
	}

	// 3. Sample Item
	item := map[string]interface{}{
		"device_type": "Laptop-X1",
		"serial":      "SN12345",
	}

	// 4. Expected calls
	// First: UpsertItemType should be called for "Laptop-X1"
	repo.On("UpsertItemType", ctx, mock.MatchedBy(func(it *domain.ItemType) bool {
		return it.Code == "Laptop-X1" && it.Name == "Laptop-X1"
	})).Return(&domain.ItemType{ID: 42, Code: "Laptop-X1"}, nil)

	// Second: UpsertAsset should be called with ItemTypeID=42
	repo.On("UpsertAsset", ctx, mock.MatchedBy(func(a *domain.Asset) bool {
		return a.ItemTypeID == 42 && *a.SerialNumber == "SN12345"
	})).Return(nil)

	// 5. Run it
	typeCache := worker.loadItemTypeCache(ctx)
	err := worker.ingestItem(ctx, mappings, item, typeCache)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestIngestWorker_ItemTypeResolutionByCode(t *testing.T) {
	repo := new(MockRepository)
	worker := NewIngestWorker(repo)
	ctx := context.Background()

	// 1. Initial Repo State: Existing item type
	repo.On("ListItemTypes", ctx, true).Return([]domain.ItemType{
		{ID: 99, Code: "SKU-99", Name: "Super Phone"},
	}, nil)

	// 2. Mappings
	mappings := []domain.IngestMapping{
		{JSONPath: "$.sku", TargetModel: domain.IngestTargetAsset, TargetField: "item_type_code"},
		{JSONPath: "$.serial", TargetModel: domain.IngestTargetAsset, TargetField: "serial_number", IsIdentity: true},
	}

	// 3. Sample Item
	item := map[string]interface{}{
		"sku":    "SKU-99",
		"serial": "SN-888",
	}

	// 4. Expected calls
	repo.On("UpsertAsset", ctx, mock.MatchedBy(func(a *domain.Asset) bool {
		return a.ItemTypeID == 99 && *a.SerialNumber == "SN-888"
	})).Return(nil)

	// 5. Run it
	typeCache := worker.loadItemTypeCache(ctx)
	err := worker.ingestItem(ctx, mappings, item, typeCache)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
