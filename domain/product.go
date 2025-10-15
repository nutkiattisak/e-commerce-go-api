package domain

import (
	"context"

	"ecommerce-go-api/entity"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error)
}

type ProductRepository interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error)
}
