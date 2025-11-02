package domain

import (
	"context"

	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error)
	GetProductsByUserID(ctx context.Context, userID uuid.UUID, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
	CreateProduct(ctx context.Context, userID uuid.UUID, req *entity.CreateProductRequest) (*entity.ProductResponse, error)
	UpdateProduct(ctx context.Context, userID uuid.UUID, productID int, req *entity.UpdateProductRequest) (*entity.ProductResponse, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, productID int) error
	ListProductsByShop(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) (*entity.ProductListResponse, error)
}

type ProductRepository interface {
	ListProducts(ctx context.Context, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	GetProductByID(ctx context.Context, productID int) (*entity.Product, error)
	ListByShopID(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) ([]*entity.Product, int64, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID int) error
	RestoreProductStock(ctx context.Context, productID int, qty int) error
}
