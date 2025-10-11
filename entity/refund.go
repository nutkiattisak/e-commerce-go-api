package entity

import (
	"time"

	"github.com/google/uuid"
)

type Refund struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_refunds_shop_order_id" json:"shopOrderId"`
	PaymentID     *uuid.UUID `gorm:"type:uuid;index:idx_refunds_payment_id" json:"paymentId"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	RefundMethod  string     `gorm:"size:100" json:"refundMethod"`
	Status        string     `gorm:"size:50;not null;default:'PENDING';index:idx_refunds_status" json:"status"`
	Reason        string     `gorm:"type:text" json:"reason"`
	BankAccount   string     `gorm:"size:100" json:"bankAccount"`
	BankName      string     `gorm:"size:100" json:"bankName"`
	TransactionID string     `gorm:"type:text" json:"transactionId"`
	RefundedAt    *time.Time `json:"refundedAt"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	ShopOrder ShopOrder `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
	Payment   *Payment  `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
}
