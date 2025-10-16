package domain

import (
	"context"

	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	GetProductByID(ctx context.Context, productID int) (*entity.Product, error)
}

type ProductRepository interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	GetProductByID(ctx context.Context, productID int) (*entity.Product, error)
	ListByShopID(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
}
