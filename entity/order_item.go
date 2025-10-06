package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopOrderID  uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_shop_order_id;uniqueIndex:uq_order_items_shop_order_product" json:"shop_order_id"`
	ProductID    uint      `gorm:"not null;uniqueIndex:uq_order_items_shop_order_product" json:"product_id"`
	Qty          int       `gorm:"not null" json:"qty"`
	UnitPrice    float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal     float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	
	// Relations
	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shop_order,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}