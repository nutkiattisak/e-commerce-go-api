package repository

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"

	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) domain.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetProvinces(ctx context.Context) ([]*entity.Province, error) {
	var out []*entity.Province
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *locationRepository) GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error) {
	var out []*entity.District
	if err := r.db.WithContext(ctx).
		Where("province_id = ? AND deleted_at IS NULL", provinceID).
		Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *locationRepository) GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error) {
	var out []*entity.SubDistrict
	if err := r.db.WithContext(ctx).
		Where("district_id = ? AND deleted_at IS NULL", districtID).
		Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
