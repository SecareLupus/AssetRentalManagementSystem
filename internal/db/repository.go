package db

import (
	"context"

	"github.com/desmond/rental-management-system/internal/domain"
)

type Repository interface {
	// ItemTypes
	GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error)
	ListItemTypes(ctx context.Context) ([]domain.ItemType, error)

	// Assets
	GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error)
	ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error)

	// RentActions
	CreateRentAction(ctx context.Context, ra *domain.RentAction) error
	GetRentActionByID(ctx context.Context, id int64) (*domain.RentAction, error)
	UpdateRentAction(ctx context.Context, ra *domain.RentAction) error
}
