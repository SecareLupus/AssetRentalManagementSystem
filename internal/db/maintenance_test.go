package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSqlRepository_RecallAssetsByItemType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSqlRepository(db)
	ctx := context.Background()

	itemTypeID := int64(1)

	mock.ExpectExec("UPDATE assets SET status = 'recalled'").
		WithArgs(sqlmock.AnyArg(), itemTypeID).
		WillReturnResult(sqlmock.NewResult(0, 5))

	err = repo.RecallAssetsByItemType(ctx, itemTypeID)
	assert.NoError(t, err)
}
