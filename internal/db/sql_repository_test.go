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

func TestSqlRepository_CreateRentAction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	ra := &domain.RentAction{
		RequesterRef: "user-1",
		CreatedByRef: "admin-1",
		Status:       "draft",
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(time.Hour),
		Items: []domain.RentActionItem{
			{ItemKind: "item_type", ItemID: 10, RequestedQuantity: 1},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO rent_actions").
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("INSERT INTO rent_action_items").
		WithArgs(1, "item_type", 10, 1, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
	mock.ExpectCommit()

	err = repo.CreateRentAction(ctx, ra)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra.ID)
	assert.Equal(t, int64(100), ra.Items[0].ID)
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

	// Mock overlapping reserved quantity
	mock.ExpectQuery("SELECT COALESCE(.+) FROM rent_action_items rai JOIN rent_actions ra").
		WithArgs(int64(10), startTime, endTime).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(3))

	avail, err := repo.GetAvailableQuantity(ctx, 10, startTime, endTime)
	assert.NoError(t, err)
	assert.Equal(t, 7, avail)
}
