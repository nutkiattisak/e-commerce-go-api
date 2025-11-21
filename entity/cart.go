package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"userId"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

type ProductSummary struct {
	ID       uint32  `json:"id"`
	Name     string  `json:"name,omitempty"`
	ImageURL *string `json:"imageUrl,omitempty"`
	Price    float64 `json:"price"`
	StockQty uint32  `json:"stockQty"`
}

type CartSummary struct {
	TotalItems uint32  `json:"totalItems"`
	TotalQty   uint32  `json:"totalQty"`
	Subtotal   float64 `json:"subtotal"`
}

type CartResponse struct {
	ID        uint32             `json:"id"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Summary   CartSummary        `json:"summary"`
	Items     []CartItemResponse `json:"items"`
}

type CartItemResponse struct {
	ID        uint32            `json:"id"`
	Product   ProductSummary    `json:"product"`
	Shop      *CartShopResponse `json:"shop,omitempty"`
	Qty       uint32            `json:"qty"`
	UnitPrice float64           `json:"unitPrice"`
	Subtotal  float64           `json:"subtotal"`
}

type CourierOption struct {
	CourierID uint32  `json:"courierId"`
	Name      string  `json:"name,omitempty"`
	Price     float64 `json:"price"`
}

type CartShopEstimate struct {
	ShopID   string         `json:"shopId"`
	Name     string         `json:"name"`
	ImageURL string         `json:"imageUrl"`
	Items    []CartItemShop `json:"items"`
	Subtotal float64        `json:"subtotal"`
	Courier  CourierOption  `json:"courier"`
}

type CartShippingEstimateResponse struct {
	Shop       []CartShopEstimate `json:"shop"`
	GrandTotal float64            `json:"grandTotal"`
}

type CartItemShop struct {
	CartItemID uint32  `json:"cartItemId"`
	ProductID  uint32  `json:"productId"`
	Qty        uint32  `json:"qty"`
	UnitPrice  float64 `json:"unitPrice"`
	Subtotal   float64 `json:"subtotal"`
}

type CartShopResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
}
