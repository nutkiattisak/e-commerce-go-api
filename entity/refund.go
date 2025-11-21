package entity

import (
	"time"

	"github.com/google/uuid"
)

type Refund struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShopOrderID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_refunds_shop_order_id" json:"shopOrderId"`
	PaymentID      *uuid.UUID `gorm:"type:uuid;index:idx_refunds_payment_id" json:"paymentId,omitempty"`
	Amount         float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	RefundMethodID *uint32    `json:"refundMethodId,omitempty"`
	RefundStatusID uint32     `gorm:"not null;default:1" json:"refundStatusId"`
	Reason         string     `gorm:"type:text" json:"reason,omitempty"`
	BankAccount    string     `gorm:"size:100" json:"bankAccount,omitempty"`
	BankName       string     `gorm:"size:100" json:"bankName,omitempty"`
	TransactionID  string     `gorm:"type:text" json:"transactionId,omitempty"`
	RefundedAt     *time.Time `json:"refundedAt,omitempty"`
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updatedAt"`

	RefundMethod *RefundMethod `gorm:"foreignKey:RefundMethodID;references:ID" json:"refundMethod,omitempty"`
	RefundStatus *RefundStatus `gorm:"foreignKey:RefundStatusID;references:ID" json:"refundStatus,omitempty"`
	Payment      *Payment      `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
	ShopOrder    *ShopOrder    `gorm:"foreignKey:ShopOrderID;references:ID" json:"shopOrder,omitempty"`
}

type CreateRefundRequest struct {
	ShopOrderID uuid.UUID `json:"shopOrderId" validate:"required"`
	Reason      string    `json:"reason" validate:"required,min=3,max=500"`
}

type ApproveRefundRequest struct {
	RefundID uuid.UUID `json:"refundId" validate:"required"`
}

type SubmitRefundBankAccountRequest struct {
	BankAccount string `json:"bankAccount" validate:"required,min=10,max=20"`
	BankName    string `json:"bankName" validate:"required,min=2,max=100"`
}

type RefundResponse struct {
	ID             uuid.UUID  `json:"id"`
	ShopOrderID    uuid.UUID  `json:"shopOrderId"`
	PaymentID      *uuid.UUID `json:"paymentId,omitempty"`
	Amount         float64    `json:"amount"`
	RefundMethodID *uint32    `json:"refundMethodId,omitempty"`
	RefundStatusID uint32     `json:"refundStatusId"`
	Reason         string     `json:"reason,omitempty"`
	BankAccount    string     `json:"bankAccount,omitempty"`
	BankName       string     `json:"bankName,omitempty"`
	TransactionID  string     `json:"transactionId,omitempty"`
	RefundedAt     *time.Time `json:"refundedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
