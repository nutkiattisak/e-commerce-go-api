package repository

import (
	"context"
	// "fmt"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) ListProducts(ctx context.Context, q *entity.ProductListRequest) (*entity.ProductListResponse, error) {
	var products []*entity.Product
	var total int64

	base := r.db.WithContext(ctx).Model(&entity.Product{}).Where("deleted_at IS NULL")

	if q != nil && q.SearchText != "" {
		like := "%" + q.SearchText + "%"
		base = base.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	perPage := 20
	if q != nil && q.PerPage > 0 {
		perPage = int(q.PerPage)
	}
	pageIdx := 0
	if q != nil && q.Page > 0 {
		pageIdx = int(q.Page) - 1
	}

	if err := base.Preload("Shop").Offset(pageIdx * perPage).Limit(perPage).Find(&products).Error; err != nil {
		return nil, err
	}

	var respItems []*entity.ProductResponse
	for _, p := range products {
		pr := &entity.ProductResponse{
			ID:          p.ID,
			Title:       p.Title,
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
			pr.Shop = &entity.ShopResponse{
				ID:       p.Shop.ID,
				Name:     p.Shop.Name,
				ImageURL: p.Shop.ImageURL,
			}
		}
		respItems = append(respItems, pr)
	}

	return &entity.ProductListResponse{Items: respItems, Total: total}, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entity.ProductResponse, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).Preload("Shop").First(&p, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	pr := &entity.ProductResponse{
		ID:          p.ID,
		Title:       p.Title,
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
		pr.Shop = &entity.ShopResponse{
			ID:       p.Shop.ID,
			Name:     p.Shop.Name,
			ImageURL: p.Shop.ImageURL,
		}
	}
	return pr, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID int) (*entity.ProductResponse, error) {
	return r.GetByID(ctx, productID)
}
