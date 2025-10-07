package entity

import (
	"time"
)

type Product struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	ImageURL    string     `gorm:"type:text" json:"image_url"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQty    int        `gorm:"not null;default:0" json:"stock_qty"`
	Status      string     `gorm:"size:50;default:'ACTIVE';index:idx_products_status" json:"status"`
	ShopID      int        `gorm:"not null;index:idx_products_shop_id" json:"shop_id"`
	CreatedAt   *time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"default:now()" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	
	// Relations
	Shop Shop `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
}