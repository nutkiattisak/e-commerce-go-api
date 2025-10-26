package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCancelled  PaymentStatus = "CANCELLED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type Payment struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID       uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_payments_order_id" json:"orderId"`
	TransactionID string     `gorm:"size:255;not null;uniqueIndex:uq_payments_transaction_id" json:"transactionId"`
	PaymentMethod string     `gorm:"size:50;not null" json:"paymentMethod"` // card, bank_transfer, promptpay
	Provider      string     `gorm:"size:100" json:"provider"`              // omise, 2c2p, etc.
	Status        string     `gorm:"size:100;not null;default:'PENDING';index:idx_payments_status" json:"status"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaidAt        *time.Time `json:"paidAt"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

type CreatePaymentRequest struct {
	PaymentMethod string  `json:"paymentMethod" validate:"required,oneof=card bank_transfer promptpay"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

type PaymentResponse struct {
	ID            uuid.UUID  `json:"id"`
	OrderID       uuid.UUID  `json:"orderId"`
	TransactionID string     `json:"transactionId"`
	PaymentMethod string     `json:"paymentMethod"`
	Provider      string     `json:"provider"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	PaidAt        *time.Time `json:"paidAt"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
