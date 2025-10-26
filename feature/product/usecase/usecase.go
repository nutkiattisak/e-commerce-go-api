package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type productUsecase struct {
	repo      domain.ProductRepository
	shopRepo  domain.ShopRepository
	validator *validator.Validate
}

func NewProductUsecase(r domain.ProductRepository, s domain.ShopRepository) domain.ProductUsecase {
	return &productUsecase{repo: r, shopRepo: s, validator: validator.New()}
}

func mapToProductResponse(p *entity.Product) *entity.ProductResponse {
	if p == nil {
		return nil
	}

	resp := &entity.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Price:       p.Price,
		StockQty:    p.StockQty,
		IsActive:    p.IsActive,
		ShopID:      p.ShopID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.Shop.ID != uuid.Nil {
		resp.Shop = &entity.ProductShopResponse{
			ID:       p.Shop.ID,
			Name:     p.Shop.Name,
			ImageURL: p.Shop.ImageURL,
		}
	}

	return resp
}

func (u *productUsecase) ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error) {
	items, total, err := u.repo.ListProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	productResponses := make([]*entity.ProductResponse, len(items))
	for i, p := range items {
		productResponses[i] = mapToProductResponse(p)
	}

	return &entity.ProductListResponse{
		Items: productResponses,
		Total: total,
	}, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error) {
	product, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return mapToProductResponse(product), nil
}

func (u *productUsecase) CreateProduct(ctx context.Context, userID uuid.UUID, req *entity.CreateProductRequest) (*entity.ProductResponse, error) {
	if err := u.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("shop not found for user: %w", err)
	}

	p := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		StockQty:    req.StockQty,
		ShopID:      shop.ID,
		IsActive:    true,
	}

	if err := u.repo.CreateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	p.Shop = *shop

	return mapToProductResponse(p), nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, userID uuid.UUID, productID int, req *entity.UpdateProductRequest) (*entity.ProductResponse, error) {
	if err := u.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	prod, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("shop not found for user: %w", err)
	}
	if prod.ShopID != shop.ID {
		return nil, fmt.Errorf("forbidden: user does not own this product")
	}

	prod.Name = req.Name
	prod.Description = req.Description
	prod.ImageURL = req.ImageURL
	prod.Price = req.Price
	prod.StockQty = req.StockQty

	if err := u.repo.UpdateProduct(ctx, prod); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return mapToProductResponse(prod), nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, userID uuid.UUID, productID int) error {
	prod, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("shop not found for user: %w", err)
	}
	if prod.ShopID != shop.ID {
		return fmt.Errorf("forbidden: user does not own this product")
	}

	if err := u.repo.DeleteProduct(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (u *productUsecase) ListProductsByShop(ctx context.Context, shopID uuid.UUID) ([]*entity.ProductResponse, error) {
	items, _, err := u.repo.ListByShopID(ctx, shopID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list products by shop: %w", err)
	}

	productResponses := make([]*entity.ProductResponse, len(items))
	for i, p := range items {
		productResponses[i] = mapToProductResponse(p)
	}

	return productResponses, nil
}

func (u *productUsecase) GetProductsByUserID(ctx context.Context, userID uuid.UUID) (*entity.ProductListResponse, error) {
	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("shop not found for user: %w", err)
	}
	items, total, err := u.repo.ListByShopID(ctx, shop.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list products for shop: %w", err)
	}

	productResponses := make([]*entity.ProductResponse, len(items))
	for i, p := range items {
		productResponses[i] = mapToProductResponse(p)
	}

	return &entity.ProductListResponse{
		Items: productResponses,
		Total: total,
	}, nil
}
