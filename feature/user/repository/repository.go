package repository

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Address, error) {
	var addrs []*entity.Address
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("SubDistrict.District.Province").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&addrs).Error
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

func (r *UserRepository) CreateAddress(ctx context.Context, addr *entity.Address) error {
	if err := r.db.WithContext(ctx).Create(addr).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetAddressByID(ctx context.Context, id int) (*entity.Address, error) {
	var a entity.Address
	if err := r.db.WithContext(ctx).First(&a, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *UserRepository) UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) error {
	updates := map[string]interface{}{
		"name":            addr.Name,
		"line1":           addr.Line1,
		"line2":           addr.Line2,
		"sub_district_id": addr.SubDistrictID,
		"district_id":     addr.DistrictID,
		"province_id":     addr.ProvinceID,
		"zipcode":         addr.Zipcode,
		"phone_number":    addr.PhoneNumber,
		"is_default":      addr.IsDefault,
		"updated_at":      time.Now(),
	}

	res := r.db.WithContext(ctx).
		Model(&entity.Address{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", addr.ID, userID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) DeleteAddress(ctx context.Context, id int, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity.Address{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *entity.User) error {
	updates := map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"phone_number": user.PhoneNumber,
		"image_url":    user.ImageURL,
		"updated_at":   time.Now(),
	}

	res := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", user.ID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
