package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (m *MockRepository) ListItemTypes(ctx context.Context) ([]domain.ItemType, error) {
	args := m.Called(ctx)
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

func (m *MockRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	args := m.Called(ctx, itemTypeID)
	return args.Get(0).([]domain.Asset), args.Error(1)
}

func (m *MockRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRepository) UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus) error {
	args := m.Called(ctx, id, status)
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

func (m *MockRepository) UpdateRentAction(ctx context.Context, ra *domain.RentAction) error {
	args := m.Called(ctx, ra)
	return args.Error(0)
}

func TestHandler_CreateItemType(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo)

	it := domain.ItemType{Code: "TEST", Name: "Test Item", Kind: domain.ItemKindSerialized}
	body, _ := json.Marshal(it)

	repo.On("CreateItemType", mock.Anything, mock.AnythingOfType("*domain.ItemType")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/catalog/item-types", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItemType(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_CreateItemType_Invalid(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo)

	it := domain.ItemType{Code: "", Name: "Test Item", Kind: domain.ItemKindSerialized} // Empty Code
	body, _ := json.Marshal(it)

	req := httptest.NewRequest(http.MethodPost, "/v1/catalog/item-types", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItemType(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetCatalog(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo)

	items := []domain.ItemType{{ID: 1, Name: "Item 1"}}
	repo.On("ListItemTypes", mock.Anything).Return(items, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/catalog/item-types", nil)
	w := httptest.NewRecorder()

	h.GetCatalog(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []domain.ItemType
	json.NewDecoder(w.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "Item 1", response[0].Name)
}
