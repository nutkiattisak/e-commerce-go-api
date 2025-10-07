package entity

import "time"

type CartItem struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    int        `gorm:"not null;uniqueIndex:uq_cart_items_cart_product" json:"cart_id"`
	ProductID int        `gorm:"not null;uniqueIndex:uq_cart_items_cart_product" json:"product_id"`
	Qty       int        `gorm:"not null;default:1" json:"qty"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	
	Cart    Cart    `gorm:"foreignKey:CartID;references:ID" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}