package entity

import "time"

type CartItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    uint      `gorm:"uniqueIndex:idx_cart_product" json:"cart_id"`
	ProductID uint      `gorm:"uniqueIndex:idx_cart_product" json:"product_id"`
	Qty       int       `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Cart    Cart    `gorm:"foreignKey:CartID;references:ID" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}