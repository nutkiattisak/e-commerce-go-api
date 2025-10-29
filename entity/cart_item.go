package entity

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    int            `gorm:"not null;uniqueIndex:uq_cart_items_cart_product" json:"cartId"`
	ProductID int            `gorm:"not null;uniqueIndex:uq_cart_items_cart_product" json:"productId"`
	Qty       int            `gorm:"not null;default:1" json:"qty"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"default:null" json:"deletedAt"`

	Cart    Cart    `gorm:"foreignKey:CartID;references:ID" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

type CartItemRequest struct {
	ProductID int `json:"productId" validate:"required,gt=0"`
	Qty       int `json:"qty" validate:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Qty int `json:"qty" validate:"required,gt=0"`
}

type EstimateShippingRequest struct {
	CartItemIDs []int `json:"cartItemIds" validate:"required,min=1,dive,gt=0"`
}
