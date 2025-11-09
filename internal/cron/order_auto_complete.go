package cron

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/timeth"
)

type OrderAutoCompleteJob struct {
	orderRepo domain.OrderRepository
}

func NewOrderAutoCompleteJob(orderRepo domain.OrderRepository) *OrderAutoCompleteJob {
	return &OrderAutoCompleteJob{
		orderRepo: orderRepo,
	}
}

func (j *OrderAutoCompleteJob) ProcessDeliveredOrders() {
	ctx := context.Background()
	startTime := timeth.Now()

	const autoCompleteDays = 7
	deliveredOrders, err := j.orderRepo.ListDeliveredOrdersOlderThan(ctx, autoCompleteDays)
	if err != nil {
		log.Printf("[CRON] Error fetching delivered orders: %v", err)
		return
	}

	if len(deliveredOrders) == 0 {
		log.Println("[CRON] No delivered orders to auto-complete")
		return
	}

	log.Printf("[CRON] Found %d delivered orders to auto-complete", len(deliveredOrders))

	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	for _, order := range deliveredOrders {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(shopOrder *entity.ShopOrder) {
			defer wg.Done()
			defer func() { <-semaphore }()

			defer func() {
				if r := recover(); r != nil {
					log.Printf("[CRON] Panic recovered while auto-completing order %s: %v", shopOrder.ID, r)
					mu.Lock()
					errorCount++
					mu.Unlock()
				}
			}()

			processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := j.autoCompleteOrder(processCtx, shopOrder); err != nil {
				log.Printf("[CRON] Error auto-completing order %s: %v", shopOrder.ID, err)
				mu.Lock()
				errorCount++
				mu.Unlock()
				return
			}

			log.Printf("[CRON] Successfully auto-completed order %s (Order Number: %s, Delivered at: %s)",
				shopOrder.ID, shopOrder.OrderNumber, shopOrder.UpdatedAt.Format(time.RFC3339))
			mu.Lock()
			successCount++
			mu.Unlock()
		}(order)
	}

	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("[CRON] Auto-complete check completed in %v - Success: %d, Errors: %d, Total: %d",
		duration, successCount, errorCount, len(deliveredOrders))
}

func (j *OrderAutoCompleteJob) autoCompleteOrder(ctx context.Context, shopOrder *entity.ShopOrder) error {
	if err := j.orderRepo.UpdateShopOrderStatus(ctx, shopOrder.ID, entity.OrderStatusCompleted); err != nil {
		return fmt.Errorf("failed to update shop order status: %w", err)
	}

	if err := j.orderRepo.UpdateShipmentStatusByShopOrderID(ctx, shopOrder.ID, entity.ShipmentStatusDelivered); err != nil {
		log.Printf("[CRON] Warning: Failed to update shipment status for order %s: %v", shopOrder.ID, err)
	}

	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		now := timeth.Now()
		orderLog := &entity.OrderLog{
			OrderID:       shopOrder.OrderID,
			ShopOrderID:   &shopOrder.ID,
			OrderStatusID: entity.OrderStatusCompleted,
			Note:          fmt.Sprintf("Auto-completed: Customer did not confirm within %d days", 7),
			CreatedAt:     &now,
		}

		if err := j.orderRepo.CreateOrderLog(logCtx, orderLog); err != nil {
			log.Printf("[CRON] Warning: Failed to create order log for %s: %v", shopOrder.ID, err)
		}
	}()

	return nil
}
