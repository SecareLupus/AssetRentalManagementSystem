package db

import (
	"context"

	"github.com/desmond/rental-management-system/internal/domain"
)

type Repository interface {
	// ItemTypes
	CreateItemType(ctx context.Context, it *domain.ItemType) error
	GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error)
	ListItemTypes(ctx context.Context) ([]domain.ItemType, error)
	UpdateItemType(ctx context.Context, it *domain.ItemType) error
	DeleteItemType(ctx context.Context, id int64) error

	// Assets
	CreateAsset(ctx context.Context, a *domain.Asset) error
	GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error)
	ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error)
	UpdateAsset(ctx context.Context, a *domain.Asset) error
	UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus) error
	DeleteAsset(ctx context.Context, id int64) error

	// RentActions
	CreateRentAction(ctx context.Context, ra *domain.RentAction) error
	GetRentActionByID(ctx context.Context, id int64) (*domain.RentAction, error)
	UpdateRentAction(ctx context.Context, ra *domain.RentAction) error
}
