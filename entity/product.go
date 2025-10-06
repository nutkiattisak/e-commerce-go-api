package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Title       string         `gorm:"size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	Price       float64        `gorm:"type:decimal(10,2)" json:"price"`
	StockQty    int            `gorm:"not null;default:0" json:"stock_qty"`
	Status      string         `gorm:"size:255" json:"status"`
	ShopID      uint           `gorm:"not null" json:"shop_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	// Relations
	Shop Shop `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
}