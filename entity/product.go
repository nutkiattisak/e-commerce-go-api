package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	ImageURL    string     `gorm:"type:text" json:"imageUrl"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQty    int        `gorm:"not null;default:0" json:"stockQty"`
	Status      string     `gorm:"size:50;default:'ACTIVE';index:idx_products_status" json:"status"`
	ShopID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_products_shop_id" json:"shopId"`
	CreatedAt   *time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`

	// Relations
	Shop Shop `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
}
