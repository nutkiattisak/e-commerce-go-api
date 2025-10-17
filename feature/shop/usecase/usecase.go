package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type shopUsecase struct {
	shopRepo    domain.ShopRepository
	productRepo domain.ProductRepository
}

func NewShopUsecase(s domain.ShopRepository, p domain.ProductRepository) domain.ShopUsecase {
	return &shopUsecase{shopRepo: s, productRepo: p}
}

func (u *shopUsecase) GetShopByID(ctx context.Context, shopID uuid.UUID) (*entity.Shop, error) {
	return u.shopRepo.GetShopByID(ctx, shopID)
}

func (u *shopUsecase) GetShopByUserID(ctx context.Context, userID uuid.UUID) (*entity.Shop, error) {
	return u.shopRepo.GetShopByUserID(ctx, userID)
}

func (u *shopUsecase) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Product, int64, error) {
	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find shop for user: %w", err)
	}

	items, total, err := u.productRepo.ListByShopID(ctx, shop.ID, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products for shop: %w", err)
	}
	return items, total, nil
}

func (u *shopUsecase) UpdateShop(ctx context.Context, shopID uuid.UUID, userID uuid.UUID, req *entity.UpdateShopRequest) (*entity.Shop, error) {
	shop, err := u.shopRepo.GetShopByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	if shop.UserID != userID {
		return nil, fmt.Errorf("forbidden: user not owner")
	}
	if req.Name != "" {
		shop.Name = req.Name
	}
	if req.Description != "" {
		shop.Description = req.Description
	}
	if req.ImageURL != "" {
		shop.ImageURL = req.ImageURL
	}
	if req.Address != "" {
		shop.Address = req.Address
	}

	if err := u.shopRepo.UpdateShop(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return shop, nil
}

func (u *shopUsecase) ListShops(ctx context.Context, req *entity.ShopListRequest) (*entity.ShopListResponse, error) {
	shops, total, err := u.shopRepo.ListShops(ctx, req)
	if err != nil {
		return nil, err
	}

	return &entity.ShopListResponse{Items: shops, Total: total}, nil
}

func (u *shopUsecase) ActivateShop(ctx context.Context, shopID uuid.UUID) error {
	return u.shopRepo.UpdateStatus(ctx, shopID, true)
}

func (u *shopUsecase) DeactivateShop(ctx context.Context, shopID uuid.UUID) error {
	return u.shopRepo.UpdateStatus(ctx, shopID, false)
}

func (u *shopUsecase) ListProductsByShop(ctx context.Context, shopID uuid.UUID) ([]*entity.Product, error) {
	return nil, nil
}
