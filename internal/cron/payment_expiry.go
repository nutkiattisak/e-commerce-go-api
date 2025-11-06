package cron

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type PaymentExpiryJob struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewPaymentExpiryJob(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) *PaymentExpiryJob {
	return &PaymentExpiryJob{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (j *PaymentExpiryJob) ProcessExpiredPayments() {
	ctx := context.Background()
	startTime := time.Now()

	expiredPayments, err := j.orderRepo.ListExpiredPayments(ctx)
	if err != nil {
		log.Printf("[CRON] Error fetching expired payments: %v", err)
		return
	}

	if len(expiredPayments) == 0 {
		log.Println("[CRON] No expired payments found")
		return
	}

	log.Printf("[CRON] Found %d expired payments to process", len(expiredPayments))

	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	for _, payment := range expiredPayments {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p *entity.Payment) {
			defer wg.Done()
			defer func() { <-semaphore }()

			defer func() {
				if r := recover(); r != nil {
					log.Printf("[CRON] Panic recovered while processing payment %s: %v", p.ID, r)
					mu.Lock()
					errorCount++
					mu.Unlock()
				}
			}()

			processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := j.processExpiredPayment(processCtx, p); err != nil {
				log.Printf("[CRON] Error processing expired payment %s: %v", p.ID, err)
				mu.Lock()
				errorCount++
				mu.Unlock()
				return
			}

			log.Printf("[CRON] Successfully processed expired payment %s (Order: %s)", p.ID, p.OrderID)
			mu.Lock()
			successCount++
			mu.Unlock()
		}(payment)
	}

	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("[CRON] Expired payment check completed in %v - Success: %d, Errors: %d, Total: %d",
		duration, successCount, errorCount, len(expiredPayments))
}

func (j *PaymentExpiryJob) processExpiredPayment(ctx context.Context, payment *entity.Payment) error {
	order, err := j.orderRepo.GetOrderByID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if err := j.orderRepo.UpdatePaymentStatus(ctx, payment.ID, 3, nil); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var processingErrors []error

	for i := range order.ShopOrders {
		wg.Add(1)

		go func(shopOrder entity.ShopOrder) {
			defer wg.Done()

			if err := j.orderRepo.UpdateShopOrderStatus(ctx, shopOrder.ID, 6); err != nil {
				mu.Lock()
				processingErrors = append(processingErrors, fmt.Errorf("failed to cancel shop order %s: %w", shopOrder.ID, err))
				mu.Unlock()
				return
			}

			j.restoreShopOrderStock(ctx, shopOrder.OrderItems)

			go func(so entity.ShopOrder) {
				logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				now := time.Now()
				shopOrderLog := &entity.OrderLog{
					OrderID:     order.ID,
					ShopOrderID: &so.ID,
					Note:        fmt.Sprintf("Order cancelled due to payment expiry (Transaction: %s)", payment.TransactionID),
					CreatedAt:   &now,
				}
				if err := j.orderRepo.CreateOrderLog(logCtx, shopOrderLog); err != nil {
					log.Printf("[CRON] Warning: Failed to create shop order log for %s: %v", so.ID, err)
				}
			}(shopOrder)
		}(order.ShopOrders[i])
	}

	wg.Wait()

	if len(processingErrors) > 0 {
		return fmt.Errorf("encountered %d errors while processing shop orders: %v", len(processingErrors), processingErrors[0])
	}

	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		now := time.Now()
		orderLog := &entity.OrderLog{
			OrderID:   order.ID,
			Note:      fmt.Sprintf("Payment expired (Transaction: %s, Expired at: %s)", payment.TransactionID, payment.ExpiresAt.Format(time.RFC3339)),
			CreatedAt: &now,
		}
		if err := j.orderRepo.CreateOrderLog(logCtx, orderLog); err != nil {
			log.Printf("[CRON] Warning: Failed to create main order log for %s: %v", order.ID, err)
		}
	}()

	return nil
}

func (j *PaymentExpiryJob) restoreShopOrderStock(ctx context.Context, orderItems []entity.OrderItem) {
	if len(orderItems) == 0 {
		return
	}

	var wg sync.WaitGroup
	for i := range orderItems {
		wg.Add(1)

		go func(item entity.OrderItem) {
			defer wg.Done()

			restoreCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := j.productRepo.RestoreProductStock(restoreCtx, item.ProductID, item.Qty); err != nil {
				log.Printf("[CRON] Warning: Failed to restore stock for product %d (qty: %d): %v", item.ProductID, item.Qty, err)
			} else {
				log.Printf("[CRON] Successfully restored stock for product %d (qty: %d)", item.ProductID, item.Qty)
			}
		}(orderItems[i])
	}

	wg.Wait()
}
