package repository

import (
	"context"

	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) domain.CourierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) ListAll(ctx context.Context) ([]*entity.Courier, error) {
	var couriers []*entity.Courier

	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&couriers).Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}
