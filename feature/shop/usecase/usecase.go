package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
)

type shopUsecase struct {
	shopRepo    domain.ShopRepository
	productRepo domain.ProductRepository
}

func NewShopUsecase(s domain.ShopRepository, p domain.ProductRepository) domain.ShopUsecase {
	return &shopUsecase{shopRepo: s, productRepo: p}
}

func (u *shopUsecase) GetShopByID(ctx context.Context, shopID uuid.UUID) (*entity.ShopResponse, error) {
	shop, err := u.shopRepo.GetShopByID(ctx, shopID)
	if err != nil {
		return nil, err
	}

	return &entity.ShopResponse{
		ID:          shop.ID,
		UserID:      shop.UserID,
		Name:        shop.Name,
		Description: shop.Description,
		ImageURL:    shop.ImageURL,
		Address:     shop.Address,
		IsActive:    shop.IsActive,
	}, nil
}

func (u *shopUsecase) GetShopByUserID(ctx context.Context, userID uuid.UUID) (*entity.ShopResponse, error) {
	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entity.ShopResponse{
		ID:          shop.ID,
		UserID:      shop.UserID,
		Name:        shop.Name,
		Description: shop.Description,
		ImageURL:    shop.ImageURL,
		Address:     shop.Address,
		IsActive:    shop.IsActive,
	}, nil
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

func (u *shopUsecase) UpdateShop(ctx context.Context, shopID uuid.UUID, userID uuid.UUID, req *entity.UpdateShopRequest) (*entity.ShopResponse, error) {
	shop, err := u.shopRepo.GetShopByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, errmap.ErrNotFound) {
			return nil, errmap.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}
	if req.Name != nil {
		shop.Name = *req.Name
	}
	if req.Description != nil {
		shop.Description = *req.Description
	}
	if req.ImageURL != nil {
		shop.ImageURL = *req.ImageURL
	}
	if req.Address != nil {
		shop.Address = *req.Address
	}

	if err := u.shopRepo.UpdateShop(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return &entity.ShopResponse{
		ID:          shop.ID,
		UserID:      shop.UserID,
		Name:        shop.Name,
		Description: shop.Description,
		ImageURL:    shop.ImageURL,
		Address:     shop.Address,
		IsActive:    shop.IsActive,
	}, nil
}

func (u *shopUsecase) ListShops(ctx context.Context, req *entity.ShopListRequest) (*entity.ShopListResponse, error) {
	shops, total, err := u.shopRepo.ListShops(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*entity.ShopResponse, 0, len(shops))
	for _, shop := range shops {
		items = append(items, &entity.ShopResponse{
			ID:          shop.ID,
			UserID:      shop.UserID,
			Name:        shop.Name,
			Description: shop.Description,
			ImageURL:    shop.ImageURL,
			Address:     shop.Address,
			IsActive:    shop.IsActive,
		})
	}

	return &entity.ShopListResponse{Items: items, Total: total}, nil
}

func (u *shopUsecase) ActivateShop(ctx context.Context, shopID uuid.UUID) error {
	return u.shopRepo.UpdateStatus(ctx, shopID, true)
}

func (u *shopUsecase) DeactivateShop(ctx context.Context, shopID uuid.UUID) error {
	return u.shopRepo.UpdateStatus(ctx, shopID, false)
}
