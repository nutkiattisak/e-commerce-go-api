package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderLog struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_order_logs_order_id" json:"order_id"`
	ShopOrderID  *uuid.UUID `gorm:"type:uuid" json:"shop_order_id"`
	Status       string     `gorm:"size:100;not null" json:"status"`
	Note         string     `gorm:"type:text" json:"note"`
	CreatedBy    *uint      `json:"created_by"`
	CreatedAt    time.Time  `gorm:"default:now();index:idx_order_logs_created_at" json:"created_at"`
	
	// Relations
	Order     Order      `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	ShopOrder *ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shop_order,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
}