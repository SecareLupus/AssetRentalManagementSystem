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

func (m *MockRepository) UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus, placeID *int64, location *string, metadata json.RawMessage) error {
	args := m.Called(ctx, id, status, placeID, location, metadata)
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

func (m *MockRepository) CreateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	args := m.Called(ctx, rr)
	return args.Error(0)
}

func (m *MockRepository) GetRentalReservationByID(ctx context.Context, id int64) (*domain.RentalReservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RentalReservation), args.Error(1)
}

func (m *MockRepository) ListRentalReservations(ctx context.Context) ([]domain.RentalReservation, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.RentalReservation), args.Error(1)
}

func (m *MockRepository) UpdateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	args := m.Called(ctx, rr)
	return args.Error(0)
}

func (m *MockRepository) UpdateRentalReservationStatus(ctx context.Context, id int64, status domain.RentalReservationStatus) error {
	args := m.Called(ctx, id, status)
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

func TestHandler_ApproveRentalReservation(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	rr := &domain.RentalReservation{
		ID:                1,
		ReservationStatus: domain.ReservationStatusPending,
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		Demands: []domain.Demand{
			{ItemKind: "item_type", ItemID: 10, Quantity: 1},
		},
	}

	repo.On("GetRentalReservationByID", mock.Anything, int64(1)).Return(rr, nil)
	repo.On("GetAvailableQuantity", mock.Anything, int64(10), mock.Anything, mock.Anything).Return(5, nil)
	repo.On("UpdateRentalReservationStatus", mock.Anything, int64(1), domain.ReservationStatusConfirmed).Return(nil)
	repo.On("AppendEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/logistics/reservations/1/approve", nil)
	w := httptest.NewRecorder()

	h.ApproveRentalReservation(w, req)

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

func (m *MockRepository) CreatePerson(ctx context.Context, p *domain.Person) error { return nil }
func (m *MockRepository) GetPerson(ctx context.Context, id int64) (*domain.Person, error) {
	return nil, nil
}
func (m *MockRepository) ListPeople(ctx context.Context) ([]domain.Person, error) {
	return nil, nil
}
func (m *MockRepository) UpdatePerson(ctx context.Context, p *domain.Person) error { return nil }
func (m *MockRepository) DeletePerson(ctx context.Context, id int64) error         { return nil }

func (m *MockRepository) CreateOrganizationRole(ctx context.Context, or *domain.OrganizationRole) error {
	return nil
}
func (m *MockRepository) ListOrganizationRoles(ctx context.Context, orgID, personID *int64) ([]domain.OrganizationRole, error) {
	return nil, nil
}
func (m *MockRepository) DeleteOrganizationRole(ctx context.Context, id int64) error { return nil }

func (m *MockRepository) CreatePlace(ctx context.Context, p *domain.Place) error { return nil }
func (m *MockRepository) GetPlace(ctx context.Context, id int64) (*domain.Place, error) {
	return nil, nil
}
func (m *MockRepository) ListPlaces(ctx context.Context, ownerID, parentID *int64) ([]domain.Place, error) {
	return nil, nil
}
func (m *MockRepository) UpdatePlace(ctx context.Context, p *domain.Place) error { return nil }
func (m *MockRepository) DeletePlace(ctx context.Context, id int64) error        { return nil }

func (m *MockRepository) CreateEvent(ctx context.Context, e *domain.Event) error { return nil }
func (m *MockRepository) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) ListEvents(ctx context.Context, companyID *int64) ([]domain.Event, error) {
	return nil, nil
}
func (m *MockRepository) UpdateEvent(ctx context.Context, e *domain.Event) error { return nil }
func (m *MockRepository) CreateDemand(ctx context.Context, d *domain.Demand) error {
	return nil
}
func (m *MockRepository) ListDemandsByReservation(ctx context.Context, reservationID int64) ([]domain.Demand, error) {
	return nil, nil
}
func (m *MockRepository) ListDemandsByEvent(ctx context.Context, eventID int64) ([]domain.Demand, error) {
	return nil, nil
}
func (m *MockRepository) UpdateDemand(ctx context.Context, d *domain.Demand) error {
	return nil
}
func (m *MockRepository) DeleteDemand(ctx context.Context, id int64) error {
	return nil
}

func (m *MockRepository) CreateCheckOutAction(ctx context.Context, coa *domain.CheckOutAction) error {
	return nil
}
func (m *MockRepository) CreateReturnAction(ctx context.Context, ra *domain.ReturnAction) error {
	return nil
}
func (m *MockRepository) ListCheckOutActions(ctx context.Context, reservationID int64) ([]domain.CheckOutAction, error) {
	return nil, nil
}
func (m *MockRepository) ListReturnActions(ctx context.Context, reservationID int64) ([]domain.ReturnAction, error) {
	return nil, nil
}
func (m *MockRepository) CreateInspection(ctx context.Context, ins *domain.InspectionSubmission) error {
	return nil
}
func (m *MockRepository) ListInspections(ctx context.Context, assetID *int64) ([]domain.InspectionSubmission, error) {
	return nil, nil
}

func (m *MockRepository) DeleteCompany(ctx context.Context, id int64) error { return nil }
func (m *MockRepository) DeleteEvent(ctx context.Context, id int64) error   { return nil }

func (m *MockRepository) GetRentalFulfillmentStatus(ctx context.Context, reservationID int64) (*domain.RentalFulfillmentStatus, error) {
	args := m.Called(ctx, reservationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RentalFulfillmentStatus), args.Error(1)
}

func (m *MockRepository) BatchCheckOut(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64, toLocationID *int64) error {
	args := m.Called(ctx, reservationID, assetIDs, agentID, toLocationID)
	return args.Error(0)
}

func (m *MockRepository) BatchReturn(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64) error {
	args := m.Called(ctx, reservationID, assetIDs, agentID)
	return args.Error(0)
}
