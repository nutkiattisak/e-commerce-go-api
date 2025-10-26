package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;index:idx_orders_user_id" json:"userId"`
	AddressID           int        `json:"addressId"`
	GrandTotal          float64    `gorm:"type:decimal(10,2);not null" json:"grandTotal"`
	CancelReason        string     `gorm:"type:text" json:"cancelReason"`
	ShippingName        string     `gorm:"size:255;not null" json:"shippingName"`
	ShippingPhone       string     `gorm:"size:15;not null" json:"shippingPhone"`
	ShippingLine1       string     `gorm:"type:text;not null" json:"shippingLine1"`
	ShippingLine2       string     `gorm:"type:text" json:"shippingLine2"`
	ShippingSubDistrict string     `gorm:"size:100;not null" json:"shippingSubDistrict"`
	ShippingDistrict    string     `gorm:"size:100;not null" json:"shippingDistrict"`
	ShippingProvince    string     `gorm:"size:100;not null" json:"shippingProvince"`
	ShippingZipcode     string     `gorm:"size:5;not null" json:"shippingZipcode"`
	PaymentMethod       string     `gorm:"size:50;not null" json:"paymentMethod"`
	CreatedAt           time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CancelledAt         *time.Time `json:"cancelledAt"`

	User       User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Address    *Address    `gorm:"foreignKey:AddressID;references:ID" json:"address,omitempty"`
	ShopOrders []ShopOrder `gorm:"foreignKey:OrderID" json:"shopOrders,omitempty"`
}

type CreateOrderRequest struct {
	CartItemIDs   []int  `json:"cartItemIds,omitempty"`
	AddressID     int    `json:"addressId,omitempty"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=card cod bank_transfer promptpay"`
}

type OrderProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	ImageURL    *string `json:"imageUrl,omitempty"`
}

type OrderItemResponse struct {
	ID        int                  `json:"id"`
	Qty       int                  `json:"qty"`
	UnitPrice float64              `json:"unitPrice"`
	Subtotal  float64              `json:"subtotal"`
	Product   OrderProductResponse `json:"product"`
}

type OrderShopResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
}

type ShopOrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	OrderID     uuid.UUID           `json:"orderId"`
	OrderNumber string              `json:"orderNumber"`
	Status      string              `json:"status"`
	Subtotal    float64             `json:"subtotal"`
	Shipping    float64             `json:"shipping"`
	GrandTotal  float64             `json:"grandTotal"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Shop        OrderShopResponse   `json:"shop"`
	OrderItems  []OrderItemResponse `json:"orderItems"`
}

type OrderResponse struct {
	ID                  uuid.UUID           `json:"id"`
	GrandTotal          float64             `json:"grandTotal"`
	CancelReason        string              `json:"cancelReason"`
	ShippingName        string              `json:"shippingName"`
	ShippingPhone       string              `json:"shippingPhone"`
	ShippingLine1       string              `json:"shippingLine1"`
	ShippingLine2       string              `json:"shippingLine2"`
	ShippingSubDistrict string              `json:"shippingSubDistrict"`
	ShippingDistrict    string              `json:"shippingDistrict"`
	ShippingProvince    string              `json:"shippingProvince"`
	ShippingZipcode     string              `json:"shippingZipcode"`
	PaymentMethod       string              `json:"paymentMethod"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	CancelledAt         *time.Time          `json:"cancelledAt"`
	ShopOrders          []ShopOrderResponse `json:"shopOrders"`
}

type OrderListResponse struct {
	ID                  uuid.UUID           `json:"id"`
	OrderID             uuid.UUID           `json:"orderId"`
	OrderNumber         string              `json:"orderNumber"`
	Status              string              `json:"status"`
	Shipping            float64             `json:"shipping"`
	GrandTotal          float64             `json:"grandTotal"`
	CancelReason        string              `json:"cancelReason"`
	ShippingName        string              `json:"shippingName"`
	ShippingPhone       string              `json:"shippingPhone"`
	ShippingLine1       string              `json:"shippingLine1"`
	ShippingLine2       string              `json:"shippingLine2"`
	ShippingSubDistrict string              `json:"shippingSubDistrict"`
	ShippingDistrict    string              `json:"shippingDistrict"`
	ShippingProvince    string              `json:"shippingProvince"`
	ShippingZipcode     string              `json:"shippingZipcode"`
	PaymentMethod       string              `json:"paymentMethod"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	CancelledAt         *time.Time          `json:"cancelledAt"`
	Shop                OrderShopResponse   `json:"shop"`
	OrderItems          []OrderItemResponse `json:"orderItems"`
}
