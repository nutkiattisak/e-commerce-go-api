package entity

import (
	"time"

	"github.com/google/uuid"
)

type ShopOrderStatus string

const (
	ShopOrderStatusPending    ShopOrderStatus = "PENDING"
	ShopOrderStatusConfirmed  ShopOrderStatus = "CONFIRMED"
	ShopOrderStatusProcessing ShopOrderStatus = "PROCESSING"
	ShopOrderStatusShipped    ShopOrderStatus = "SHIPPED"
	ShopOrderStatusDelivered  ShopOrderStatus = "DELIVERED"
	ShopOrderStatusCancelled  ShopOrderStatus = "CANCELLED"
	ShopOrderStatusRefunded   ShopOrderStatus = "REFUNDED"
)

type ShopOrder struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index:idx_shop_orders_order_id;uniqueIndex:uq_shop_orders_order_shop" json:"order_id"`
	ShopID      int       `gorm:"not null;index:idx_shop_orders_shop_id;index:idx_shop_orders_shop_status;index:idx_shop_orders_shop_created;uniqueIndex:uq_shop_orders_order_shop" json:"shop_id"`
	Status      string    `gorm:"size:100;not null;default:'PENDING';index:idx_shop_orders_status;index:idx_shop_orders_shop_status" json:"status"`
	Subtotal    float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Shipping    float64   `gorm:"type:decimal(10,2);not null" json:"shipping"`
	GrandTotal  float64   `gorm:"type:decimal(10,2);not null" json:"grand_total"`
	CreatedAt   time.Time `gorm:"not null;default:now();index:idx_shop_orders_shop_created" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
	
	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	Shop  Shop  `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
}