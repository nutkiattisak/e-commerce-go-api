package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name" validate:"required,min=3,max=255"`
	Description string         `gorm:"type:text" json:"description" validate:"omitempty,max=2000"`
	ImageURL    *string        `gorm:"type:text" json:"imageUrl,omitempty" validate:"omitempty,url"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
	StockQty    int            `gorm:"not null;default:0" json:"stockQty" validate:"gte=0"`
	IsActive    bool           `gorm:"default:true;index:idx_products_is_active" json:"isActive"`
	ShopID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_products_shop_id" json:"shopId"`
	CreatedAt   *time.Time     `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   *time.Time     `gorm:"default:now()" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Shop Shop `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
}

type ShopResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ImageURL string    `json:"imageUrl,omitempty"`
}

type ProductResponse struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ImageURL    *string       `json:"imageUrl,omitempty"`
	Price       float64       `json:"price"`
	StockQty    int           `json:"stockQty"`
	IsActive    bool          `json:"isActive"`
	ShopID      uuid.UUID     `json:"shopId"`
	CreatedAt   *time.Time    `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time    `json:"updatedAt,omitempty"`
	Shop        *ShopResponse `json:"shop,omitempty"`
}

type ProductListRequest struct {
	Page       uint32 `query:"page" validate:"omitempty,min=1"`
	PerPage    uint32 `query:"perPage" validate:"omitempty,min=1,max=100"`
	SearchText string `query:"searchText"`
}

type ProductListResponse struct {
	Items []*ProductResponse `json:"items"`
	Total int64              `json:"total"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,max=2000"`
	ImageURL    *string `json:"imageUrl,omitempty" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	StockQty    int     `json:"stockQty" validate:"gte=0"`
}
