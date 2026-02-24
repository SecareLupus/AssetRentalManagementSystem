package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateScheduledDelivery(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	sd := domain.ScheduledDelivery{EventID: 1, TargetDate: time.Now(), Notes: "Test"}
	body, _ := json.Marshal(sd)

	repo.On("CreateScheduledDelivery", mock.Anything, mock.AnythingOfType("*domain.ScheduledDelivery")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/logistics/deliveries", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateScheduledDelivery(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_CreateShipment(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	eventID := int64(1)
	s := domain.Shipment{ScheduledDeliveryID: &eventID, Status: "Preparing"}
	body, _ := json.Marshal(s)

	repo.On("CreateShipment", mock.Anything, mock.AnythingOfType("*domain.Shipment")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/logistics/shipments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateShipment(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestHandler_AllocateAssets(t *testing.T) {
	repo := new(MockRepository)
	h := NewHandler(repo, nil)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"asset_ids": []int64{101, 102},
	})

	repo.On("AllocateAssetsToShipment", mock.Anything, int64(42), []int64{101, 102}, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/logistics/shipments/42/allocate", bytes.NewBuffer(reqBody))
	// Mock user ID in context (JWT unmarshals numbers as float64)
	claims := map[string]interface{}{
		"user_id": float64(1),
	}
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.AllocateAssets(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
