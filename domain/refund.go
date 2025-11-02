package domain

import (
	"context"

	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type RefundUsecase interface {
	CreateRefund(ctx context.Context, userID uuid.UUID, req entity.CreateRefundRequest) (*entity.RefundResponse, error)
	ApproveRefund(ctx context.Context, userID uuid.UUID, refundID uuid.UUID) (*entity.RefundResponse, error)
	SubmitRefundBankAccount(ctx context.Context, userID uuid.UUID, refundID uuid.UUID, req entity.SubmitRefundBankAccountRequest) (*entity.RefundResponse, error)
}

type RefundRepository interface {
	CreateRefund(ctx context.Context, refund *entity.Refund) error
	GetRefundByID(ctx context.Context, id uuid.UUID) (*entity.Refund, error)
	UpdateRefundStatus(ctx context.Context, id uuid.UUID, statusID int) error
	UpdateRefundBankAccount(ctx context.Context, id uuid.UUID, bankAccount, bankName string) error
}
