package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderLog struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_order_logs_order_id" json:"orderId"`
	ShopOrderID *uuid.UUID `gorm:"type:uuid" json:"shopOrderId"`
	Status      string     `gorm:"size:100;not null" json:"status"`
	Note        string     `gorm:"type:text" json:"note"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"createdBy"`
	CreatedAt   *time.Time `gorm:"default:now();index:idx_order_logs_created_at" json:"createdAt"`
	
	
	Order     Order      `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	ShopOrder *ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
}