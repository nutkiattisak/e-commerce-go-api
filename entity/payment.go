package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_payments_order_id" json:"orderId"`
	PaymentMethodID uint32     `gorm:"size:50;not null" json:"paymentMethodId"` // 1, 2, 3
	PaymentStatusID uint32     `gorm:"size:100;not null;default:1;index:idx_payments_status" json:"paymentStatusId"`
	TransactionID   string     `gorm:"size:255;not null;uniqueIndex:uq_payments_transaction_id" json:"transactionId"`
	Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaidAt          *time.Time `json:"paidAt"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	CreatedAt       time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

type CreatePaymentRequest struct {
	PaymentMethodID uint32  `json:"paymentMethodId" validate:"required,oneof=1 2 3 4"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
}

type PaymentResponse struct {
	ID              uuid.UUID  `json:"id"`
	OrderID         uuid.UUID  `json:"orderId"`
	TransactionID   string     `json:"transactionId"`
	PaymentMethodID uint32     `json:"paymentMethodId"`
	PaymentStatusID uint32     `json:"paymentStatusId"`
	Amount          float64    `json:"amount"`
	PaidAt          *time.Time `json:"paidAt"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
