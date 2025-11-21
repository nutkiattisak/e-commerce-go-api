package domain

import (
	"context"

	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	GetProductByID(ctx context.Context, productID uint32) (*entity.ProductResponse, error)
	GetProductsByUserID(ctx context.Context, userID uuid.UUID, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	CreateProduct(ctx context.Context, userID uuid.UUID, req *entity.CreateProductRequest) (*entity.ProductResponse, error)
	UpdateProduct(ctx context.Context, userID uuid.UUID, productID uint32, req *entity.UpdateProductRequest) (*entity.ProductResponse, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, productID uint32) error
	ListProductsByShop(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
}

type ProductRepository interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	GetProductByID(ctx context.Context, productID uint32) (*entity.Product, error)
	ListByShopID(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID uint32) error
	RestoreProductStock(ctx context.Context, productID uint32, qty uint32) error
}
