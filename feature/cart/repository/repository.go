package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var c entity.Cart
	if err := r.db.WithContext(ctx).First(&c, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *cartRepository) EnsureCartForUser(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	c, err := r.GetCartByUserID(ctx, userID)
	if err == nil {
		return c, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cart := &entity.Cart{UserID: userID}
	if createErr := r.db.WithContext(ctx).Create(cart).Error; createErr != nil {
		return nil, createErr
	}
	return cart, nil
}

func (r *cartRepository) ListCartItems(ctx context.Context, cartID uint32) ([]*entity.CartItem, error) {
	var items []*entity.CartItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Shop").
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepository) AddCartItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *cartRepository) UpsertCartItem(ctx context.Context, item *entity.CartItem) (*entity.CartItem, bool, error) {
	var existing entity.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ? AND deleted_at IS NULL", item.CartID, item.ProductID).
		First(&existing).Error

	if err == nil {
		existing.Qty += item.Qty
		if updateErr := r.db.WithContext(ctx).Save(&existing).Error; updateErr != nil {
			return nil, false, updateErr
		}

		var out entity.CartItem
		if err := r.db.WithContext(ctx).
			Preload("Product").
			Preload("Product.Shop").
			Where("id = ?", existing.ID).
			First(&out).Error; err != nil {
			return nil, false, err
		}

		return &out, false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	r.db.WithContext(ctx).Unscoped().
		Where("cart_id = ? AND product_id = ? AND deleted_at IS NOT NULL", item.CartID, item.ProductID).
		Delete(&entity.CartItem{})

	if createErr := r.db.WithContext(ctx).Create(item).Error; createErr != nil {
		return nil, false, createErr
	}

	var out entity.CartItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Shop").
		Where("id = ?", item.ID).
		First(&out).Error; err != nil {
		return nil, false, err
	}

	return &out, true, nil
}

func (r *cartRepository) GetCartItemByID(ctx context.Context, id uint32) (*entity.CartItem, error) {
	var it entity.CartItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Shop").
		Preload("Cart").
		First(&it, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *cartRepository) GetCartItemByCartAndProduct(ctx context.Context, cartID uint32, productID uint32) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	if err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ? AND deleted_at IS NULL", cartID, productID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartRepository) GetCartItemByUserAndProduct(ctx context.Context, userID uuid.UUID, productID uint32) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	if err := r.db.WithContext(ctx).
		Table("cart_items as ci").
		Select("ci.*").
		Joins("join carts c on c.id = ci.cart_id").
		Where("c.user_id = ? AND ci.product_id = ? AND ci.deleted_at IS NULL", userID, productID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartRepository) GetCartItemsByIDs(ctx context.Context, id []uint32) ([]*entity.CartItem, error) {
	var cartItems []*entity.CartItem
	if len(id) == 0 {
		return cartItems, nil
	}

	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Shop").
		Where("id IN ? AND deleted_at IS NULL", id).
		Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (r *cartRepository) UpdateCartItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *cartRepository) DeleteCartItem(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.CartItem{}).Error
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID uint32) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error
}
