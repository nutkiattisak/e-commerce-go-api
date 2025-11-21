package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_user_id" json:"userId"`
	AddressID           uint32    `json:"addressId"`
	GrandTotal          float64   `gorm:"type:decimal(10,2);not null" json:"grandTotal"`
	ShippingName        string    `gorm:"size:255;not null" json:"shippingName"`
	ShippingPhone       string    `gorm:"size:15;not null" json:"shippingPhone"`
	ShippingLine1       string    `gorm:"type:text;not null" json:"shippingLine1"`
	ShippingLine2       string    `gorm:"type:text" json:"shippingLine2"`
	ShippingSubDistrict string    `gorm:"size:100;not null" json:"shippingSubDistrict"`
	ShippingDistrict    string    `gorm:"size:100;not null" json:"shippingDistrict"`
	ShippingProvince    string    `gorm:"size:100;not null" json:"shippingProvince"`
	ShippingZipcode     string    `gorm:"size:5;not null" json:"shippingZipcode"`
	PaymentMethodID     uint32    `json:"paymentMethodId"`
	CreatedAt           time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	User          User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Address       *Address       `gorm:"foreignKey:AddressID;references:ID" json:"address,omitempty"`
	ShopOrders    []ShopOrder    `gorm:"foreignKey:OrderID" json:"shopOrders,omitempty"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID" json:"paymentMethod,omitempty"`
}

type CreateOrderRequest struct {
	CartItemIDs     []uint32 `json:"cartItemIds" validate:"required"`
	AddressID       uint32   `json:"addressId" validate:"required,gt=0"`
	PaymentMethodID uint32   `json:"paymentMethodId" validate:"required,gt=0"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=500"`
}

type UpdateOrderStatusRequest struct {
	OrderStatusID *uint32 `json:"orderStatusId" validate:"required,oneof=2 4 5"`
}

type AddItemToCartRequest struct {
	ProductID uint32 `json:"productId" validate:"required,gt=0" example:"1"`
	Qty       uint32 `json:"qty" validate:"required,gt=0" example:"2"`
}

type OrderProductResponse struct {
	ID          uint32  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	ImageURL    *string `json:"imageUrl,omitempty"`
}

type OrderItemResponse struct {
	ID        uint32               `json:"id"`
	Qty       uint32               `json:"qty"`
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
	OrderStatusID uint32              `json:"orderStatusId"`
	Subtotal      float64             `json:"subtotal"`
	Shipping      float64             `json:"shipping"`
	GrandTotal    float64             `json:"grandTotal"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	Shop          OrderShopResponse   `json:"shop"`
	OrderItems    []OrderItemResponse `json:"orderItems"`
	Timeline      []OrderTimelineItem `json:"timeline,omitempty"`
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
	PaymentMethodID     uint32              `json:"paymentMethodId"`
	ShopOrders          []ShopOrderResponse `json:"shopOrders"`
	Timeline            []OrderTimelineItem `json:"timeline,omitempty"`
}

type OrderListResponse struct {
	ID                  uuid.UUID           `json:"id"`
	OrderID             uuid.UUID           `json:"orderId"`
	OrderNumber         string              `json:"orderNumber"`
	OrderStatusID       uint32              `json:"orderStatusId"`
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
	PaymentMethodID     uint32              `json:"paymentMethodId"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	Shop                OrderShopResponse   `json:"shop"`
	OrderItems          []OrderItemResponse `json:"orderItems"`
	Timeline            []OrderTimelineItem `json:"timeline,omitempty"`
}

type ShopOrderListResponse struct {
	ID                  uuid.UUID           `json:"id"`
	OrderID             uuid.UUID           `json:"orderId"`
	OrderNumber         string              `json:"orderNumber"`
	OrderStatusID       uint32              `json:"orderStatusId"`
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
	PaymentMethodID     uint32              `json:"paymentMethodId"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	OrderItems          []OrderItemResponse `json:"orderItems"`
	Timeline            []OrderTimelineItem `json:"timeline,omitempty"`
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
	OrderStatusID *uint32 `query:"orderStatusId"`
}
