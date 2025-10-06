package entity

import (
	"time"

	"github.com/google/uuid"
)

type Refund struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_refunds_shop_order_id" json:"shop_order_id"`
	PaymentID     *uuid.UUID `gorm:"type:uuid;index:idx_refunds_payment_id" json:"payment_id"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	RefundMethod  string     `gorm:"size:100" json:"refund_method"`
	Status        string     `gorm:"size:50;not null;default:'PENDING';index:idx_refunds_status" json:"status"`
	Reason        string     `gorm:"type:text" json:"reason"`
	BankAccount   string     `gorm:"size:100" json:"bank_account"`
	BankName      string     `gorm:"size:100" json:"bank_name"`
	TransactionID string     `gorm:"type:text" json:"transaction_id"`
	RefundedAt    *time.Time `json:"refunded_at"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	
	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shop_order,omitempty"`
	Payment   *Payment  `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
}