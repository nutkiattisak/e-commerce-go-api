package entity

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_shipments_shop_order" json:"shop_order_id"`
	CourierID   int        `gorm:"not null" json:"courier_id"`
	TrackingNo  string     `gorm:"size:100;not null;uniqueIndex:uq_shipments_tracking_no" json:"tracking_no"`
	Status      string     `gorm:"size:100;default:'PENDING';index:idx_shipments_status" json:"status"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	
	// Relations
	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shop_order,omitempty"`
	Courier   Courier   `gorm:"foreignKey:CourierID;references:ID" json:"courier,omitempty"`
}