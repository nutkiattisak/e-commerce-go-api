package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID       uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_payments_order_id" json:"orderId"`
	Provider      string     `gorm:"size:100" json:"provider"`
	Status        string     `gorm:"size:100;not null;default:'PENDING';index:idx_payments_status" json:"status"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	TransactionID string     `gorm:"type:text" json:"transactionId"`
	PaidAt        *time.Time `json:"paidAt"`
	CreatedAt     *time.Time `gorm:"default:now()" json:"createdAt"`

	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}
