package domain

import (
	"context"

	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type ShopUsecase interface {
	GetShopByID(ctx context.Context, shopID uuid.UUID) (*entity.Shop, error)
	GetShopByUserID(ctx context.Context, userID uuid.UUID) (*entity.Shop, error)
	UpdateShop(ctx context.Context, shopID uuid.UUID, userID uuid.UUID, req *entity.UpdateShopRequest) (*entity.Shop, error)
	ListShops(ctx context.Context, req *entity.ShopListRequest) (*entity.ShopListResponse, error)
	ActivateShop(ctx context.Context, shopID uuid.UUID) error
	DeactivateShop(ctx context.Context, shopID uuid.UUID) error
}

type ShopRepository interface {
	GetShopByID(ctx context.Context, id uuid.UUID) (*entity.Shop, error)
	GetShopByUserID(ctx context.Context, userID uuid.UUID) (*entity.Shop, error)
	GetShopsByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Shop, error)
	UpdateShop(ctx context.Context, shop *entity.Shop) error
	ListShops(ctx context.Context, req *entity.ShopListRequest) ([]*entity.Shop, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	ListShopCouriersByShopIDs(ctx context.Context, shopIDs []uuid.UUID) ([]*entity.ShopCourier, error)
}
