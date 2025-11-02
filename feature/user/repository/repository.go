package repository

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var u entity.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Address, error) {
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

func (r *userRepository) CreateAddress(ctx context.Context, addr *entity.Address) error {
	if err := r.db.WithContext(ctx).Create(addr).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAddressByID(ctx context.Context, id int) (*entity.Address, error) {
	var address entity.Address

	err := r.db.WithContext(ctx).
		Preload("SubDistrict", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name_th, name_en, district_id").
				Preload("District", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, name_th, name_en, province_id").
						Preload("Province", func(db *gorm.DB) *gorm.DB {
							return db.Select("id, name_th, name_en")
						})
				})
		}).
		Preload("District", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name_th, name_en, province_id").
				Preload("Province", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, name_th, name_en")
				})
		}).
		Preload("Province", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name_th, name_en")
		}).
		First(&address, "id = ? AND deleted_at IS NULL", id).Error

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *userRepository) UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) error {
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

func (r *userRepository) DeleteAddress(ctx context.Context, id int, userID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND is_default = false", id, userID).
		Delete(&entity.Address{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return errmap.ErrForbidden
}

func (r *userRepository) UpdateProfile(ctx context.Context, user *entity.User) error {
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
