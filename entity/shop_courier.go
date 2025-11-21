package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopCourier struct {
	ID        uint32         `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_shop_couriers_shop_id" json:"shopId"`
	CourierID uint32         `gorm:"not null;index:idx_shop_couriers_courier_id" json:"courierId"`
	Rate      float64        `gorm:"type:decimal(10,2)" json:"rate,omitempty"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Shop    *Shop    `gorm:"foreignKey:ShopID;references:ID" json:"shop,omitempty"`
	Courier *Courier `gorm:"foreignKey:CourierID;references:ID" json:"courier,omitempty"`
}
