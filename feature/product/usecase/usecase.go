package usecase

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type productUsecase struct {
	repo     domain.ProductRepository
	shopRepo domain.ShopRepository
}

func NewProductUsecase(r domain.ProductRepository, s domain.ShopRepository) domain.ProductUsecase {
	return &productUsecase{repo: r, shopRepo: s}
}

func (u *productUsecase) ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error) {
	return u.repo.ListProducts(ctx, q)
}

func (u *productUsecase) GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error) {
	return u.repo.GetProductByID(ctx, productID)
}
