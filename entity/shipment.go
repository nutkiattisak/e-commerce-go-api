package entity

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID      uuid.UUID  `gorm:"type:uuid;not null" json:"shopOrderId"`
	CourierID        uint32     `gorm:"not null" json:"courierId"`
	TrackingNo       string     `gorm:"size:100;not null" json:"trackingNo"`
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	ShippedAt        *time.Time `json:"shippedAt"`
	DeliveredAt      *time.Time `json:"deliveredAt"`
	ShipmentStatusID uint32     `gorm:"not null;default:1" json:"shipmentStatusId"`
	UpdatedAt        time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	ShopOrder      ShopOrder       `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
	Courier        Courier         `gorm:"foreignKey:CourierID;references:ID" json:"courier,omitempty"`
	ShipmentStatus *ShipmentStatus `gorm:"foreignKey:ShipmentStatusID;references:ID" json:"shipmentStatus,omitempty"`
}

type AddShipmentRequest struct {
	CourierID  uint32 `json:"courierId" validate:"required,gt=0"`
	TrackingNo string `json:"trackingNo" validate:"required,min=3,max=100"`
}

type ShipmentResponse struct {
	ID               uuid.UUID               `json:"id"`
	ShopOrderID      uuid.UUID               `json:"shopOrderId"`
	CourierID        uint32                  `json:"courierId"`
	Courier          *CourierListResponse    `json:"courier,omitempty"`
	TrackingNo       string                  `json:"trackingNo"`
	ShipmentStatusID uint32                  `json:"shipmentStatusId"`
	ShipmentStatus   *ShipmentStatusResponse `json:"shipmentStatus,omitempty"`
	CreatedAt        time.Time               `json:"createdAt"`
	ShippedAt        *time.Time              `json:"shippedAt"`
}
