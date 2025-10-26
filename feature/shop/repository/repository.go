package repository

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) domain.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) GetShopByID(ctx context.Context, id uuid.UUID) (*entity.Shop, error) {
	var s entity.Shop
	if err := r.db.WithContext(ctx).First(&s, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *shopRepository) GetShopByUserID(ctx context.Context, userID uuid.UUID) (*entity.Shop, error) {
	var s entity.Shop
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *shopRepository) GetShopsByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Shop, error) {
	var shops []*entity.Shop
	if len(ids) == 0 {
		return shops, nil
	}
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *shopRepository) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *shopRepository) UpdateShop(ctx context.Context, shop *entity.Shop) error {
	return r.db.WithContext(ctx).Save(shop).Error
}

func (r *shopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Shop{}, "id = ?", id).Error
}

func (r *shopRepository) ListShops(ctx context.Context, req *entity.ShopListRequest) ([]*entity.Shop, int64, error) {
	var shops []*entity.Shop
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Shop{}).Where("deleted_at IS NULL")
	if req.SearchText != "" {
		q = q.Where("name ILIKE ?", "%"+req.SearchText+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.PerPage == 0 {
		req.PerPage = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}
	offset := (req.Page - 1) * req.PerPage
	if err := q.Offset(offset).Limit(req.PerPage).Find(&shops).Error; err != nil {
		return nil, 0, err
	}
	return shops, total, nil
}

func (r *shopRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).Model(&entity.Shop{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *shopRepository) ListShopCouriersByShopIDs(ctx context.Context, shopIDs []uuid.UUID) ([]*entity.ShopCourier, error) {
	var scs []*entity.ShopCourier
	if len(shopIDs) == 0 {
		return scs, nil
	}

	if err := r.db.WithContext(ctx).
		Preload("Courier").
		Where("shop_id IN ? AND deleted_at IS NULL", shopIDs).
		Find(&scs).Error; err != nil {
		return nil, err
	}

	return scs, nil
}
