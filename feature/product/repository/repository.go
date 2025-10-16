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

func (r *productRepository) ListProducts(ctx context.Context, q *entity.ProductListRequest) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	base := r.db.WithContext(ctx).Model(&entity.Product{}).Where("deleted_at IS NULL")

	if q != nil && q.SearchText != "" {
		like := "%" + q.SearchText + "%"
		base = base.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).Preload("Shop").First(&p, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).Preload("Shop").First(&p, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) ListByShopID(ctx context.Context, shopID uuid.UUID, q *entity.ProductListRequest) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	productQuery := r.db.WithContext(ctx).Model(&entity.Product{}).Where("deleted_at IS NULL AND shop_id = ?", shopID)

	if q != nil && q.SearchText != "" {
		like := "%" + q.SearchText + "%"
		productQuery = productQuery.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	if err := productQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	perPage := 20
	if q != nil && q.PerPage > 0 {
		perPage = int(q.PerPage)
	}

	page := 0
	if q != nil && q.Page > 0 {
		page = int(q.Page) - 1
	}

	if err := productQuery.Preload("Shop").Offset(page * perPage).Limit(perPage).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}
