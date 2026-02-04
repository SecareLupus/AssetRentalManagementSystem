package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSqlRepository_BuildSpecs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	bs := &domain.BuildSpec{
		Version: "v1.0.0",
	}

	mock.ExpectQuery("INSERT INTO build_specs").
		WithArgs("v1.0.0", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.CreateBuildSpec(ctx, bs)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), bs.ID)

	mock.ExpectQuery("SELECT id, version, hardware_config, software_config, firmware_url, metadata, created_at, updated_at FROM build_specs WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "version", "hardware_config", "software_config", "firmware_url", "metadata", "created_at", "updated_at"}).
			AddRow(1, "v1.0.0", []byte("{}"), []byte("{}"), nil, []byte("{}"), time.Now(), time.Now()))

	ret, err := repo.GetBuildSpecByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, "v1.0.0", ret.Version)
}

func TestSqlRepository_Provisioning(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	assetID := int64(100)
	buildSpecID := int64(1)
	performedBy := "tester"

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE assets SET status = 'maintenance', provisioning_status = 'flashing', current_build_spec_id = \\$1 WHERE id = \\$2").
		WithArgs(buildSpecID, assetID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("INSERT INTO provision_actions").
		WithArgs(assetID, &buildSpecID, domain.ProvisionStarted, performedBy, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(500))
	mock.ExpectCommit()

	pa, err := repo.StartProvisioning(ctx, assetID, buildSpecID, performedBy)
	assert.NoError(t, err)
	assert.NotNil(t, pa)
	assert.Equal(t, int64(500), pa.ID)

	// Complete Provisioning
	mock.ExpectBegin()
	mock.ExpectQuery("UPDATE provision_actions SET status = 'completed', notes = \\$1, completed_at = \\$2 WHERE id = \\$3 RETURNING asset_id").
		WithArgs("done", sqlmock.AnyArg(), int64(500)).
		WillReturnRows(sqlmock.NewRows([]string{"asset_id"}).AddRow(assetID))

	mock.ExpectExec("UPDATE assets SET status = 'available', provisioning_status = 'ready' WHERE id = \\$1").
		WithArgs(assetID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.CompleteProvisioning(ctx, 500, "done")
	assert.NoError(t, err)
}
