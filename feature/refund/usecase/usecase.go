package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
)

type refundUsecase struct {
	refundRepo domain.RefundRepository
	orderRepo  domain.OrderRepository
	shopRepo   domain.ShopRepository
}

func NewRefundUsecase(refundRepo domain.RefundRepository, orderRepo domain.OrderRepository, shopRepo domain.ShopRepository) domain.RefundUsecase {
	return &refundUsecase{
		refundRepo: refundRepo,
		orderRepo:  orderRepo,
		shopRepo:   shopRepo,
	}
}

func (u *refundUsecase) CreateRefund(ctx context.Context, userID uuid.UUID, req entity.CreateRefundRequest) (*entity.RefundResponse, error) {
	shopOrder, err := u.orderRepo.GetShopOrderByID(ctx, req.ShopOrderID)
	if err != nil {
		return nil, fmt.Errorf("shop order not found: %w", err)
	}

	shop, err := u.shopRepo.GetShopByID(ctx, shopOrder.ShopID)
	if err != nil {
		return nil, err
	}
	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	if shopOrder.OrderStatusID != entity.OrderStatusPending {
		return nil, fmt.Errorf("refund can only be created for orders with status = 1")
	}

	payment, err := u.orderRepo.GetPaymentByOrderID(ctx, shopOrder.OrderID)
	if err != nil {
		return nil, fmt.Errorf("payment not found for this order: %w", err)
	}

	if payment.PaymentStatusID != entity.PaymentStatusRefunded {
		return nil, fmt.Errorf("refund can only be created when payment_status = 6")
	}

	refundMethodID := entity.RefundMethodBankTransfer
	if payment.PaymentMethodID == entity.PaymentMethodCreditCard {
		refundMethodID = entity.RefundMethodCreditCard
	}

	now := time.Now()
	refund := &entity.Refund{
		ShopOrderID:    req.ShopOrderID,
		PaymentID:      &payment.ID,
		Amount:         shopOrder.GrandTotal,
		RefundMethodID: &refundMethodID,
		RefundStatusID: entity.RefundStatusPending,
		Reason:         req.Reason,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.refundRepo.CreateRefund(ctx, refund); err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return &entity.RefundResponse{
		ID:             refund.ID,
		ShopOrderID:    refund.ShopOrderID,
		PaymentID:      refund.PaymentID,
		Amount:         refund.Amount,
		RefundMethodID: refund.RefundMethodID,
		RefundStatusID: refund.RefundStatusID,
		Reason:         refund.Reason,
		CreatedAt:      refund.CreatedAt,
		UpdatedAt:      refund.UpdatedAt,
	}, nil
}

func (u *refundUsecase) ApproveRefund(ctx context.Context, userID uuid.UUID, refundID uuid.UUID) (*entity.RefundResponse, error) {
	refund, err := u.refundRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}

	shopOrder, err := u.orderRepo.GetShopOrderByID(ctx, refund.ShopOrderID)
	if err != nil {
		return nil, err
	}

	shop, err := u.shopRepo.GetShopByID(ctx, shopOrder.ShopID)
	if err != nil {
		return nil, err
	}
	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	if refund.RefundStatusID != entity.RefundStatusPending {
		return nil, fmt.Errorf("refund is not in pending status")
	}

	if err := u.refundRepo.UpdateRefundStatus(ctx, refundID, entity.RefundStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to approve refund: %w", err)
	}

	updatedRefund, err := u.refundRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, err
	}

	return &entity.RefundResponse{
		ID:             updatedRefund.ID,
		ShopOrderID:    updatedRefund.ShopOrderID,
		PaymentID:      updatedRefund.PaymentID,
		Amount:         updatedRefund.Amount,
		RefundMethodID: updatedRefund.RefundMethodID,
		RefundStatusID: updatedRefund.RefundStatusID,
		Reason:         updatedRefund.Reason,
		RefundedAt:     updatedRefund.RefundedAt,
		CreatedAt:      updatedRefund.CreatedAt,
		UpdatedAt:      updatedRefund.UpdatedAt,
	}, nil
}

func (u *refundUsecase) SubmitRefundBankAccount(ctx context.Context, userID uuid.UUID, refundID uuid.UUID, req entity.SubmitRefundBankAccountRequest) (*entity.RefundResponse, error) {
	refund, err := u.refundRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}

	shopOrder, err := u.orderRepo.GetShopOrderByID(ctx, refund.ShopOrderID)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.GetOrderByID(ctx, shopOrder.OrderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	if refund.RefundStatusID != entity.RefundStatusPending {
		return nil, fmt.Errorf("cannot submit bank account for non-pending refund")
	}

	if err := u.refundRepo.UpdateRefundBankAccount(ctx, refundID, req.BankAccount, req.BankName); err != nil {
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	updatedRefund, err := u.refundRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, err
	}

	return &entity.RefundResponse{
		ID:             updatedRefund.ID,
		ShopOrderID:    updatedRefund.ShopOrderID,
		PaymentID:      updatedRefund.PaymentID,
		Amount:         updatedRefund.Amount,
		RefundMethodID: updatedRefund.RefundMethodID,
		RefundStatusID: updatedRefund.RefundStatusID,
		Reason:         updatedRefund.Reason,
		BankAccount:    updatedRefund.BankAccount,
		BankName:       updatedRefund.BankName,
		CreatedAt:      updatedRefund.CreatedAt,
		UpdatedAt:      updatedRefund.UpdatedAt,
	}, nil
}
