package entity

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_shipments_shop_order" json:"shopOrderId"`
	CourierID   int        `gorm:"not null" json:"courierId"`
	TrackingNo  string     `gorm:"size:100;not null;uniqueIndex:uq_shipments_tracking_no" json:"trackingNo"`
	Status      string     `gorm:"size:100;default:'PENDING';index:idx_shipments_status" json:"status"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	ShippedAt   *time.Time `json:"shippedAt"`
	DeliveredAt *time.Time `json:"deliveredAt"`

	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
	Courier   Courier   `gorm:"foreignKey:CourierID;references:ID" json:"courier,omitempty"`
}
