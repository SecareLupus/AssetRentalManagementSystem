package api

import (
	"bytes"
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
