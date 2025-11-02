package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type refundRepository struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) domain.RefundRepository {
	return &refundRepository{db: db}
}

func (r *refundRepository) CreateRefund(ctx context.Context, refund *entity.Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

func (r *refundRepository) GetRefundByID(ctx context.Context, id uuid.UUID) (*entity.Refund, error) {
	var refund entity.Refund
	err := r.db.WithContext(ctx).
		Preload("ShopOrder").
		Preload("Payment").
		First(&refund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *refundRepository) UpdateRefundStatus(ctx context.Context, id uuid.UUID, statusID int) error {
	updates := map[string]interface{}{
		"refund_status_id": statusID,
		"updated_at":       time.Now(),
	}

	if statusID == entity.RefundStatusCompleted {
		now := time.Now()
		updates["refunded_at"] = now
	}

	return r.db.WithContext(ctx).Model(&entity.Refund{}).Where("id = ?", id).Updates(updates).Error
}

func (r *refundRepository) UpdateRefundBankAccount(ctx context.Context, id uuid.UUID, bankAccount, bankName string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Refund{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"bank_account": bankAccount,
			"bank_name":    bankName,
			"updated_at":   time.Now(),
		}).Error
}
