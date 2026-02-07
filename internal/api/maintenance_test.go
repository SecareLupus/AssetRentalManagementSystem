package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Re-using MockRepository from handlers_test.go would be better, but they are in same package.
// Wait, they ARE in same package (api). So I should only define it once.
// Let me check if they are in the same package.

func TestHandler_RecallItemTypeAssets(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	repo.On("RecallAssetsByItemType", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/fleet/item-types/1/recall", nil)
	w := httptest.NewRecorder()

	h.RecallItemTypeAssets(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_RepairAsset(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	repo.On("UpdateAssetStatus", mock.Anything, int64(1), domain.AssetStatusMaintenance, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/inventory/assets/1/repair", nil)
	w := httptest.NewRecorder()

	h.RepairAsset(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_RefurbishAsset(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	asset := &domain.Asset{ID: 1}
	repo.On("GetAssetByID", mock.Anything, int64(1)).Return(asset, nil)
	repo.On("UpdateAsset", mock.Anything, mock.MatchedBy(func(a *domain.Asset) bool {
		return a.Status == domain.AssetStatusMaintenance && *a.CurrentBuildSpecID == int64(10)
	})).Return(nil)

	body, _ := json.Marshal(struct {
		BuildSpecID int64 `json:"build_spec_id"`
	}{BuildSpecID: 10})
	req := httptest.NewRequest(http.MethodPost, "/v1/inventory/assets/1/refurbish", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.RefurbishAsset(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
