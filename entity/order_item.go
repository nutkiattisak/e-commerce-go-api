package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopOrderID uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_shop_order_id;uniqueIndex:uq_order_items_shop_order_product" json:"shopOrderId"`
	ProductID   int       `gorm:"not null;uniqueIndex:uq_order_items_shop_order_product" json:"productId"`
	Qty         int       `gorm:"not null" json:"qty"`
	UnitPrice   float64   `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	Subtotal    float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`

	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}
