package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of db.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateItemType(ctx context.Context, it *domain.ItemType) error {
	args := m.Called(ctx, it)
	return args.Error(0)
}

func (m *MockRepository) GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ItemType), args.Error(1)
}

func (m *MockRepository) ListItemTypes(ctx context.Context, includeInactive bool) ([]domain.ItemType, error) {
	args := m.Called(ctx, includeInactive)
	return args.Get(0).([]domain.ItemType), args.Error(1)
}

func (m *MockRepository) UpdateItemType(ctx context.Context, it *domain.ItemType) error {
	args := m.Called(ctx, it)
	return args.Error(0)
}

func (m *MockRepository) DeleteItemType(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateAsset(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRepository) GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Asset), args.Error(1)
}

func (m *MockRepository) ListAssets(ctx context.Context) ([]domain.Asset, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Asset), args.Error(1)
}

func (m *MockRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	args := m.Called(ctx, itemTypeID)
	return args.Get(0).([]domain.Asset), args.Error(1)
}

func (m *MockRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRepository) UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus, location *string, metadata json.RawMessage) error {
	args := m.Called(ctx, id, status, location, metadata)
	return args.Error(0)
}

func (m *MockRepository) RecallAssetsByItemType(ctx context.Context, itemTypeID int64) error {
	args := m.Called(ctx, itemTypeID)
	return args.Error(0)
}

func (m *MockRepository) DeleteAsset(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateRentAction(ctx context.Context, ra *domain.RentAction) error {
	args := m.Called(ctx, ra)
	return args.Error(0)
}

func (m *MockRepository) GetRentActionByID(ctx context.Context, id int64) (*domain.RentAction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RentAction), args.Error(1)
}

func (m *MockRepository) ListRentActions(ctx context.Context) ([]domain.RentAction, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.RentAction), args.Error(1)
}

func (m *MockRepository) UpdateRentAction(ctx context.Context, ra *domain.RentAction) error {
	args := m.Called(ctx, ra)
	return args.Error(0)
}

func (m *MockRepository) UpdateRentActionStatus(ctx context.Context, id int64, status domain.RentActionStatus, timestampField string, timestampValue time.Time) error {
	args := m.Called(ctx, id, status, timestampField, timestampValue)
	return args.Error(0)
}

func (m *MockRepository) GetAvailableQuantity(ctx context.Context, itemTypeID int64, startTime, endTime time.Time) (int, error) {
	args := m.Called(ctx, itemTypeID, startTime, endTime)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) AddMaintenanceLog(ctx context.Context, ml *domain.MaintenanceLog) error {
	args := m.Called(ctx, ml)
	return args.Error(0)
}

func (m *MockRepository) ListMaintenanceLogs(ctx context.Context, assetID int64) ([]domain.MaintenanceLog, error) {
	args := m.Called(ctx, assetID)
	return args.Get(0).([]domain.MaintenanceLog), args.Error(1)
}

func (m *MockRepository) CreateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	args := m.Called(ctx, it)
	return args.Error(0)
}

func (m *MockRepository) UpdateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	args := m.Called(ctx, it)
	return args.Error(0)
}

func (m *MockRepository) DeleteInspectionTemplate(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) ListInspectionTemplates(ctx context.Context) ([]domain.InspectionTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InspectionTemplate), args.Error(1)
}

func (m *MockRepository) GetInspectionTemplate(ctx context.Context, id int64) (*domain.InspectionTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InspectionTemplate), args.Error(1)
}

func (m *MockRepository) GetInspectionTemplatesForItemType(ctx context.Context, itemTypeID int64) ([]domain.InspectionTemplate, error) {
	args := m.Called(ctx, itemTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InspectionTemplate), args.Error(1)
}

func (m *MockRepository) SetItemTypeInspections(ctx context.Context, itemTypeID int64, templateIDs []int64) error {
	args := m.Called(ctx, itemTypeID, templateIDs)
	return args.Error(0)
}

func (m *MockRepository) CreateInspectionSubmission(ctx context.Context, is *domain.InspectionSubmission) error {
	args := m.Called(ctx, is)
	return args.Error(0)
}

func (m *MockRepository) CreateBuildSpec(ctx context.Context, bs *domain.BuildSpec) error {
	args := m.Called(ctx, bs)
	return args.Error(0)
}

func (m *MockRepository) GetBuildSpecByID(ctx context.Context, id int64) (*domain.BuildSpec, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BuildSpec), args.Error(1)
}

func (m *MockRepository) ListBuildSpecs(ctx context.Context) ([]domain.BuildSpec, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.BuildSpec), args.Error(1)
}

func (m *MockRepository) StartProvisioning(ctx context.Context, assetID int64, buildSpecID int64, performedBy string) (*domain.ProvisionAction, error) {
	args := m.Called(ctx, assetID, buildSpecID, performedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProvisionAction), args.Error(1)
}

func (m *MockRepository) CompleteProvisioning(ctx context.Context, actionID int64, notes string) error {
	args := m.Called(ctx, actionID, notes)
	return args.Error(0)
}

func TestHandler_CreateItemType(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	it := domain.ItemType{Code: "TEST", Name: "Test Item", Kind: domain.ItemKindSerialized}
	body, _ := json.Marshal(it)

	repo.On("CreateItemType", mock.Anything, mock.AnythingOfType("*domain.ItemType")).Return(nil)
	repo.On("AppendEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil) // Added

	req := httptest.NewRequest(http.MethodPost, "/v1/catalog/item-types", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItemType(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_CreateItemType_Invalid(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	it := domain.ItemType{Code: "", Name: "Test Item", Kind: domain.ItemKindSerialized} // Empty Code
	body, _ := json.Marshal(it)

	req := httptest.NewRequest(http.MethodPost, "/v1/catalog/item-types", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItemType(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetCatalog(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	items := []domain.ItemType{{ID: 1, Name: "Item 1"}}
	repo.On("ListItemTypes", mock.Anything, false).Return(items, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/catalog/item-types", nil)
	w := httptest.NewRecorder()

	h.GetCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []domain.ItemType
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "Item 1", response[0].Name)
}

func TestHandler_ApproveRentAction(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	ra := &domain.RentAction{
		ID:        1,
		Status:    domain.RentActionStatusPending,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Items: []domain.RentActionItem{
			{ItemKind: "item_type", ItemID: 10, RequestedQuantity: 1},
		},
	}

	repo.On("GetRentActionByID", mock.Anything, int64(1)).Return(ra, nil)
	repo.On("GetAvailableQuantity", mock.Anything, int64(10), mock.Anything, mock.Anything).Return(5, nil)
	repo.On("UpdateRentActionStatus", mock.Anything, int64(1), domain.RentActionStatusApproved, "approved_at", mock.Anything).Return(nil)
	repo.On("AppendEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil) // Added

	req := httptest.NewRequest(http.MethodPost, "/v1/rent-actions/1/approve", nil)
	w := httptest.NewRecorder()

	h.ApproveRentAction(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func (m *MockRepository) CreateUser(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) AppendEvent(ctx context.Context, tx *sql.Tx, event *domain.OutboxEvent) error {

	args := m.Called(ctx, tx, event)
	return args.Error(0)
}

func (m *MockRepository) GetPendingEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]domain.OutboxEvent), args.Error(1)
}

func (m *MockRepository) MarkEventProcessed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) MarkEventFailed(ctx context.Context, id int64, errMessage string) error {
	args := m.Called(ctx, id, errMessage)
	return args.Error(0)
}

func (m *MockRepository) GetAvailabilityTimeline(ctx context.Context, itemTypeID int64, start, end time.Time) ([]domain.AvailabilityPoint, error) {
	args := m.Called(ctx, itemTypeID, start, end)
	return args.Get(0).([]domain.AvailabilityPoint), args.Error(1)
}

func (m *MockRepository) GetShortageAlerts(ctx context.Context) ([]domain.ShortageAlert, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ShortageAlert), args.Error(1)
}

func (m *MockRepository) GetMaintenanceForecast(ctx context.Context) ([]domain.MaintenanceForecast, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.MaintenanceForecast), args.Error(1)
}

func (m *MockRepository) ListWebhooks(ctx context.Context) ([]domain.WebhookConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WebhookConfig), args.Error(1)
}
func (m *MockRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DashboardStats), args.Error(1)
}
func (m *MockRepository) BulkRecallAssets(ctx context.Context, ids []int64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockRepository) CreateCompany(ctx context.Context, c *domain.Company) error { return nil }
func (m *MockRepository) GetCompany(ctx context.Context, id int64) (*domain.Company, error) {
	return nil, nil
}
func (m *MockRepository) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	return nil, nil
}
func (m *MockRepository) UpdateCompany(ctx context.Context, c *domain.Company) error { return nil }
func (m *MockRepository) CreateContact(ctx context.Context, c *domain.Contact) error { return nil }
func (m *MockRepository) GetContact(ctx context.Context, id int64) (*domain.Contact, error) {
	return nil, nil
}
func (m *MockRepository) ListContacts(ctx context.Context, companyID *int64) ([]domain.Contact, error) {
	return nil, nil
}
func (m *MockRepository) UpdateContact(ctx context.Context, c *domain.Contact) error { return nil }
func (m *MockRepository) CreateSite(ctx context.Context, s *domain.Site) error       { return nil }
func (m *MockRepository) GetSite(ctx context.Context, id int64) (*domain.Site, error) {
	return nil, nil
}
func (m *MockRepository) ListSites(ctx context.Context, companyID *int64) ([]domain.Site, error) {
	return nil, nil
}
func (m *MockRepository) UpdateSite(ctx context.Context, s *domain.Site) error         { return nil }
func (m *MockRepository) CreateLocation(ctx context.Context, l *domain.Location) error { return nil }
func (m *MockRepository) GetLocation(ctx context.Context, id int64) (*domain.Location, error) {
	return nil, nil
}
func (m *MockRepository) ListLocations(ctx context.Context, siteID, parentID *int64) ([]domain.Location, error) {
	return nil, nil
}
func (m *MockRepository) UpdateLocation(ctx context.Context, l *domain.Location) error { return nil }
func (m *MockRepository) CreateEvent(ctx context.Context, e *domain.Event) error       { return nil }
func (m *MockRepository) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) ListEvents(ctx context.Context, companyID *int64) ([]domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) UpdateEvent(ctx context.Context, e *domain.Event) error { return nil }
func (m *MockRepository) CreateEventAssetNeed(ctx context.Context, ean *domain.EventAssetNeed) error {
	return nil
}
func (m *MockRepository) ListEventAssetNeeds(ctx context.Context, eventID int64) ([]domain.EventAssetNeed, error) {
	return nil, nil
}
func (m *MockRepository) UpdateEventAssetNeed(ctx context.Context, ean *domain.EventAssetNeed) error {
	return nil
}
func (m *MockRepository) CreateInspection(ctx context.Context, ins *domain.InspectionSubmission) error {
	return nil
}
func (m *MockRepository) ListInspections(ctx context.Context, assetID *int64) ([]domain.InspectionSubmission, error) {
	return nil, nil
}
