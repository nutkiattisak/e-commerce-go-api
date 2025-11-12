package repository

import (
	"context"

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

	base := r.db.WithContext(ctx).Model(&entity.Product{}).Where("is_active = true AND deleted_at IS NULL")

	if q != nil && q.SearchText != "" {
		like := "%" + q.SearchText + "%"
		base = base.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := q.Page
	if page == 0 {
		page = 1
	}
	perPage := q.PerPage
	if perPage == 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	if err := base.Preload("Shop").Offset(int(offset)).Limit(int(perPage)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).Preload("Shop").First(&p, "id = ? AND is_active = true AND deleted_at IS NULL", id).Error; err != nil {
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

func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	updates := map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"image_url":   product.ImageURL,
		"price":       product.Price,
		"stock_qty":   product.StockQty,
		"is_active":   product.IsActive,
		"updated_at":  product.UpdatedAt,
	}
	res := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ? AND deleted_at IS NULL", product.ID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, productID int) error {
	res := r.db.WithContext(ctx).
		Where("id = ?", productID).
		Delete(&entity.Product{})
	return res.Error
}

func (r *productRepository) RestoreProductStock(ctx context.Context, productID int, qty int) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("stock_qty", gorm.Expr("stock_qty + ?", qty))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
