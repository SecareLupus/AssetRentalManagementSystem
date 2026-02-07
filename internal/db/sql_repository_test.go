package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSqlRepository_GetItemTypeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "code", "name", "kind", "is_active", "supported_features", "created_by_user_id", "updated_by_user_id", "schema_org", "metadata", "created_at", "updated_at"}).
		AddRow(1, "SKU123", "Item A", "serialized", true, []byte("{}"), nil, nil, []byte("{}"), []byte("{}"), time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM item_types WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	it, err := repo.GetItemTypeByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, it)
	assert.Equal(t, int64(1), it.ID)
	assert.Equal(t, "SKU123", it.Code)
}

func TestSqlRepository_CreateRentalReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	rr := &domain.RentalReservation{
		ID:                0,
		ReservationName:   "Test Reservation",
		ReservationStatus: domain.ReservationStatusPending,
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		Demands: []domain.Demand{
			{ItemKind: "item_type", ItemID: 10, Quantity: 1},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO rental_reservations").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("INSERT INTO demands").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
	mock.ExpectCommit()

	err = repo.CreateRentalReservation(ctx, rr)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rr.ID)
	assert.Equal(t, int64(100), rr.Demands[0].ID)
}

func TestSqlRepository_GetAvailableQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)

	// Mock total assets
	mock.ExpectQuery("SELECT COUNT(.+) FROM assets WHERE item_type_id = \\$1 AND status != 'retired'").
		WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	// Mock overlapping reserved quantity (Confirmed status)
	mock.ExpectQuery("SELECT COALESCE(.+) FROM demands d JOIN rental_reservations rr").
		WithArgs(int64(10), startTime, endTime).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(3))

	// Mock ad-hoc usage
	mock.ExpectQuery("SELECT COUNT(.+) FROM assets WHERE item_type_id = \\$1 AND status IN").
		WithArgs(int64(10), startTime).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	avail, err := repo.GetAvailableQuantity(ctx, 10, startTime, endTime)
	assert.NoError(t, err)
	assert.Equal(t, 7, avail)
}
