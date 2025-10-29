package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_user_id" json:"userId"`
	AddressID           int       `json:"addressId"`
	GrandTotal          float64   `gorm:"type:decimal(10,2);not null" json:"grandTotal"`
	ShippingName        string    `gorm:"size:255;not null" json:"shippingName"`
	ShippingPhone       string    `gorm:"size:15;not null" json:"shippingPhone"`
	ShippingLine1       string    `gorm:"type:text;not null" json:"shippingLine1"`
	ShippingLine2       string    `gorm:"type:text" json:"shippingLine2"`
	ShippingSubDistrict string    `gorm:"size:100;not null" json:"shippingSubDistrict"`
	ShippingDistrict    string    `gorm:"size:100;not null" json:"shippingDistrict"`
	ShippingProvince    string    `gorm:"size:100;not null" json:"shippingProvince"`
	ShippingZipcode     string    `gorm:"size:5;not null" json:"shippingZipcode"`
	PaymentMethodID     int       `json:"paymentMethodId"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	User          User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Address       *Address       `gorm:"foreignKey:AddressID;references:ID" json:"address,omitempty"`
	ShopOrders    []ShopOrder    `gorm:"foreignKey:OrderID" json:"shopOrders,omitempty"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID" json:"paymentMethod,omitempty"`
}

type CreateOrderRequest struct {
	CartItemIDs     []int `json:"cartItemIds" validate:"required"`
	AddressID       int   `json:"addressId" validate:"required,gt=0"`
	PaymentMethodID int   `json:"paymentMethodId" validate:"required,oneof=1 2 3 4"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=500"`
}

type UpdateOrderStatusRequest struct {
	OrderStatusID *int `json:"orderStatusId" validate:"required,oneof=1 2 3 4 5 6"`
}

type AddItemToCartRequest struct {
	ProductID int `json:"productId" validate:"required,gt=0" example:"1"`
	Qty       int `json:"qty" validate:"required,gt=0" example:"2"`
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
	ID            uuid.UUID           `json:"id"`
	OrderID       uuid.UUID           `json:"orderId"`
	OrderNumber   string              `json:"orderNumber"`
	OrderStatusID int                 `json:"orderStatusId"`
	Subtotal      float64             `json:"subtotal"`
	Shipping      float64             `json:"shipping"`
	GrandTotal    float64             `json:"grandTotal"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	Shop          OrderShopResponse   `json:"shop"`
	OrderItems    []OrderItemResponse `json:"orderItems"`
}

type OrderResponse struct {
	ID                  uuid.UUID           `json:"id"`
	GrandTotal          float64             `json:"grandTotal"`
	ShippingName        string              `json:"shippingName"`
	ShippingPhone       string              `json:"shippingPhone"`
	ShippingLine1       string              `json:"shippingLine1"`
	ShippingLine2       string              `json:"shippingLine2"`
	ShippingSubDistrict string              `json:"shippingSubDistrict"`
	ShippingDistrict    string              `json:"shippingDistrict"`
	ShippingProvince    string              `json:"shippingProvince"`
	ShippingZipcode     string              `json:"shippingZipcode"`
	PaymentMethodID     int                 `json:"paymentMethodId"`
	ShopOrders          []ShopOrderResponse `json:"shopOrders"`
}

type OrderListResponse struct {
	ID                  uuid.UUID           `json:"id"`
	OrderID             uuid.UUID           `json:"orderId"`
	OrderNumber         string              `json:"orderNumber"`
	OrderStatusID       int                 `json:"orderStatusId"`
	Shipping            float64             `json:"shipping"`
	GrandTotal          float64             `json:"grandTotal"`
	ShippingName        string              `json:"shippingName"`
	ShippingPhone       string              `json:"shippingPhone"`
	ShippingLine1       string              `json:"shippingLine1"`
	ShippingLine2       string              `json:"shippingLine2"`
	ShippingSubDistrict string              `json:"shippingSubDistrict"`
	ShippingDistrict    string              `json:"shippingDistrict"`
	ShippingProvince    string              `json:"shippingProvince"`
	ShippingZipcode     string              `json:"shippingZipcode"`
	PaymentMethodID     int                 `json:"paymentMethodId"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	Shop                OrderShopResponse   `json:"shop"`
	OrderItems          []OrderItemResponse `json:"orderItems"`
}

type ShopOrderListResponse struct {
	ID                  uuid.UUID           `json:"id"`
	OrderID             uuid.UUID           `json:"orderId"`
	OrderNumber         string              `json:"orderNumber"`
	OrderStatusID       int                 `json:"orderStatusId"`
	Shipping            float64             `json:"shipping"`
	GrandTotal          float64             `json:"grandTotal"`
	ShippingName        string              `json:"shippingName"`
	ShippingPhone       string              `json:"shippingPhone"`
	ShippingLine1       string              `json:"shippingLine1"`
	ShippingLine2       string              `json:"shippingLine2"`
	ShippingSubDistrict string              `json:"shippingSubDistrict"`
	ShippingDistrict    string              `json:"shippingDistrict"`
	ShippingProvince    string              `json:"shippingProvince"`
	ShippingZipcode     string              `json:"shippingZipcode"`
	PaymentMethodID     int                 `json:"paymentMethodId"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	OrderItems          []OrderItemResponse `json:"orderItems"`
}

type OrderListPaginationResponse struct {
	Items []*OrderListResponse `json:"items"`
	Total int64                `json:"total"`
}

type ShopOrderListPaginationResponse struct {
	Items []*ShopOrderListResponse `json:"items"`
	Total int64                    `json:"total"`
}

type OrderGroupListPaginationResponse struct {
	Items []*OrderResponse `json:"items"`
	Total int64            `json:"total"`
}

type OrderListRequest struct {
	Page          uint64  `query:"page" validate:"omitempty,min=1" example:"1"`
	PerPage       uint64  `query:"perPage" validate:"omitempty,min=1,max=100" example:"10"`
	SearchText    *string `query:"searchText" example:""`
	OrderStatusID *int    `query:"orderStatusId"`
}
