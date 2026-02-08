package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSqlRepository_CreateAsset_Lifecycle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	t.Run("Default Location and Available Status", func(t *testing.T) {
		a := &domain.Asset{
			ItemTypeID: 1,
			AssetTag:   stringPtr("TAG123"),
		}

		// 1. GetDefaultInternalPlace - Query existing
		mock.ExpectQuery("SELECT (.+) FROM places WHERE is_internal = TRUE LIMIT 1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "contained_in_place_id", "owner_id", "category", "address", "is_internal", "presumed_demands", "metadata", "created_at", "updated_at"}).
				AddRow(101, "Main Warehouse", nil, nil, nil, "site", []byte("{}"), true, []byte("{}"), []byte("{}"), time.Now(), time.Now()))

		// 2. GetPlace (for status inference)
		mock.ExpectQuery("SELECT (.+) FROM places WHERE id = \\$1").
			WithArgs(int64(101)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "contained_in_place_id", "owner_id", "category", "address", "is_internal", "presumed_demands", "metadata", "created_at", "updated_at"}).
				AddRow(101, "Main Warehouse", nil, nil, nil, "site", []byte("{}"), true, []byte("{}"), []byte("{}"), time.Now(), time.Now()))

		// 3. Insert Asset
		mock.ExpectQuery("INSERT INTO assets").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateAsset(ctx, a)
		assert.NoError(t, err)
		assert.Equal(t, int64(101), *a.PlaceID)
		assert.Equal(t, domain.AssetStatusAvailable, a.Status)
	})

	t.Run("External Location and Deployed Status", func(t *testing.T) {
		a := &domain.Asset{
			ItemTypeID: 1,
			AssetTag:   stringPtr("TAG456"),
			PlaceID:    int64Ptr(202),
		}

		// 1. GetPlace (for status inference)
		mock.ExpectQuery("SELECT (.+) FROM places WHERE id = \\$1").
			WithArgs(int64(202)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "contained_in_place_id", "owner_id", "category", "address", "is_internal", "presumed_demands", "metadata", "created_at", "updated_at"}).
				AddRow(202, "Client Site", nil, nil, nil, "site", []byte("{}"), false, []byte("{}"), []byte("{}"), time.Now(), time.Now()))

		// 2. Insert Asset
		mock.ExpectQuery("INSERT INTO assets").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		err := repo.CreateAsset(ctx, a)
		assert.NoError(t, err)
		assert.Equal(t, domain.AssetStatusDeployed, a.Status)
	})

	t.Run("Component Tracking", func(t *testing.T) {
		a := &domain.Asset{
			ItemTypeID: 1,
			AssetTag:   stringPtr("TAG789"),
			PlaceID:    int64Ptr(101),
			Components: []domain.Component{
				{Name: "CPU", SerialNumber: "CPU123"},
			},
		}

		// 1. GetPlace (for status inference)
		mock.ExpectQuery("SELECT (.+) FROM places WHERE id = \\$1").
			WithArgs(int64(101)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "contained_in_place_id", "owner_id", "category", "address", "is_internal", "presumed_demands", "metadata", "created_at", "updated_at"}).
				AddRow(101, "Main Warehouse", nil, nil, nil, "site", []byte("{}"), true, []byte("{}"), []byte("{}"), time.Now(), time.Now()))

		// 2. Insert Asset (Metadata should contain components)
		mock.ExpectQuery("INSERT INTO assets").
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

		err := repo.CreateAsset(ctx, a)
		assert.NoError(t, err)
	})
}

func stringPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64    { return &i }
