package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/timeth"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var c entity.Cart
	if err := r.db.WithContext(ctx).First(&c, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *orderRepository) EnsureCartForUser(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	c, err := r.GetCartByUserID(ctx, userID)
	if err == nil {
		return c, nil
	}

	cart := &entity.Cart{UserID: userID}
	if err := r.db.WithContext(ctx).Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *orderRepository) ListCartItems(ctx context.Context, cartID int) ([]*entity.CartItem, error) {
	var items []*entity.CartItem
	if err := r.db.WithContext(ctx).
		Preload("Product", "deleted_at IS NULL").
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepository) AddCartItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *orderRepository) UpsertCartItem(ctx context.Context, item *entity.CartItem) (*entity.CartItem, error) {
	var result *entity.CartItem

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product entity.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return err
		}

		var existing entity.CartItem
		err := tx.Where("cart_id = ? AND product_id = ? AND deleted_at IS NULL", item.CartID, item.ProductID).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if product.StockQty < item.Qty {
					return errmap.ErrInsufficientStock
				}

				if err := tx.Create(item).Error; err != nil {
					return err
				}

				result = item
				return nil
			}
			return err
		}

		newQty := existing.Qty + item.Qty
		if product.StockQty < newQty {
			return errmap.ErrInsufficientStock
		}

		existing.Qty = newQty
		existing.UpdatedAt = timeth.Now()
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		result = &existing
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *orderRepository) GetCartItemByID(ctx context.Context, id int) (*entity.CartItem, error) {
	var it entity.CartItem
	if err := r.db.WithContext(ctx).First(&it, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *orderRepository) UpdateCartItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *orderRepository) DeleteCartItem(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.CartItem{}).Error
}

func (r *orderRepository) ClearCart(ctx context.Context, cartID int) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) CreateShopOrder(ctx context.Context, so *entity.ShopOrder) error {
	return r.db.WithContext(ctx).Create(so).Error
}

func (r *orderRepository) CreateOrderItems(ctx context.Context, items []*entity.OrderItem) error {
	for _, it := range items {
		if err := r.db.WithContext(ctx).Create(it).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepository) CreateFullOrder(ctx context.Context, order *entity.Order, shopOrders []*entity.ShopOrder, orderItemsByShop map[string][]*entity.OrderItem, payment *entity.Payment, cartID int, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := timeth.Now()

		order.CreatedAt = now
		order.UpdatedAt = now
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, so := range shopOrders {
			so.OrderID = order.ID
			so.CreatedAt = now
			so.UpdatedAt = now
			if err := tx.Create(so).Error; err != nil {
				return err
			}

			items, ok := orderItemsByShop[so.ShopID.String()]
			if ok {
				for _, it := range items {
					it.ShopOrderID = so.ID

					var product entity.Product
					if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
						Where("id = ? AND is_active = ?", it.ProductID, true).
						First(&product).Error; err != nil {
						return err
					}

					if product.StockQty < it.Qty {
						return errmap.ErrInsufficientStock
					}

					if len(items) > 0 {
						if err := tx.Create(&items).Error; err != nil {
							return fmt.Errorf("failed to create order items: %w", err)
						}
					}

					res := tx.Model(&entity.Product{}).
						Where("id = ?", it.ProductID).
						Update("stock_qty", gorm.Expr("stock_qty - ?", it.Qty))

					if res.Error != nil {
						return res.Error
					}

					if res.RowsAffected == 0 {
						return errmap.ErrInsufficientStock
					}
				}
			}
		}

		if payment != nil {
			payment.OrderID = order.ID
			payment.CreatedAt = now
			payment.UpdatedAt = now
			if err := tx.Create(payment).Error; err != nil {
				return err
			}
		}

		for _, so := range shopOrders {
			shopOrderLog := &entity.OrderLog{
				OrderID:       order.ID,
				ShopOrderID:   &so.ID,
				OrderStatusID: entity.OrderStatusPending,
				CreatedBy:     &userID,
				CreatedAt:     &now,
			}
			if err := tx.Create(shopOrderLog).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) ListOrdersByUser(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var total int64

	page := req.Page
	if page == 0 {
		page = 1
	}
	perPage := req.PerPage
	if perPage == 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	query := r.db.WithContext(ctx).Model(&entity.Order{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("ShopOrders").
		Preload("ShopOrders.Shop").
		Preload("ShopOrders.OrderItems").
		Preload("ShopOrders.OrderItems.Product").
		Order("created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *orderRepository) ListShopOrdersByUserID(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) ([]*entity.ShopOrder, int64, error) {
	var shopOrders []*entity.ShopOrder
	var total int64

	page := req.Page
	if page == 0 {
		page = 1
	}
	perPage := req.PerPage
	if perPage == 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	query := r.db.WithContext(ctx).Model(&entity.ShopOrder{}).
		Joins("JOIN orders ON orders.id = shop_orders.order_id").
		Where("orders.user_id = ?", userID)

	if req.SearchText != nil {
		searchPattern := "%" + *req.SearchText + "%"
		query = query.Where("shop_orders.order_number LIKE ?", searchPattern)
	}

	if req.OrderStatusID != nil {
		query = query.Where("shop_orders.order_status_id = ?", *req.OrderStatusID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Shop").
		Preload("Order").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("shop_orders.created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&shopOrders).Error; err != nil {
		return nil, 0, err
	}

	return shopOrders, total, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Address").
		Preload("ShopOrders").
		Preload("ShopOrders.Shop").
		Preload("ShopOrders.OrderItems").
		Preload("ShopOrders.OrderItems.Product").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) ListShopOrdersByShopID(ctx context.Context, shopID uuid.UUID, req entity.OrderListRequest) ([]*entity.ShopOrder, int64, error) {
	var shopOrders []*entity.ShopOrder
	var total int64

	page := req.Page
	if page == 0 {
		page = 1
	}
	perPage := req.PerPage
	if perPage == 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	query := r.db.WithContext(ctx).Model(&entity.ShopOrder{}).Where("shop_id = ?", shopID)

	if req.SearchText != nil {
		searchPattern := "%" + *req.SearchText + "%"
		query = query.Where("order_number LIKE ? ", searchPattern)
	}

	if req.OrderStatusID != nil {
		query = query.Where("order_status_id = ?", *req.OrderStatusID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Order").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&shopOrders).Error; err != nil {
		return nil, 0, err
	}

	return shopOrders, total, nil
}

func (r *orderRepository) GetShopOrderByID(ctx context.Context, id uuid.UUID) (*entity.ShopOrder, error) {
	var so entity.ShopOrder
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Order.User").
		Preload("Order.Address").
		Preload("Shop").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&so, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &so, nil
}

func (r *orderRepository) UpdateShopOrderStatus(ctx context.Context, id uuid.UUID, OrderStatusID int) error {
	updates := map[string]interface{}{
		"order_status_id": OrderStatusID,
		"updated_at":      timeth.Now(),
	}
	return r.db.WithContext(ctx).Model(&entity.ShopOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (r *orderRepository) CancelShopOrder(ctx context.Context, id uuid.UUID, reason string) error {
	return r.db.WithContext(ctx).Model(&entity.ShopOrder{}).Where("id = ?", id).Updates(map[string]interface{}{"order_status_id": entity.OrderStatusCancelled}).Error
}

func (r *orderRepository) AddShipment(ctx context.Context, s *entity.Shipment) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *orderRepository) GetShipmentByShopOrderID(ctx context.Context, shopOrderID uuid.UUID) (*entity.Shipment, error) {
	var shipment entity.Shipment
	err := r.db.WithContext(ctx).
		Preload("Courier").
		Where("shop_order_id = ?", shopOrderID).
		First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *orderRepository) UpdateShipmentStatusByShopOrderID(ctx context.Context, shopOrderID uuid.UUID, shipmentStatusID int) error {
	return r.db.WithContext(ctx).
		Model(&entity.Shipment{}).
		Where("shop_order_id = ?", shopOrderID).
		Update("shipment_status_id", shipmentStatusID).Error
}

func (r *orderRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	var existing entity.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND payment_status_id != ?", payment.OrderID, entity.PaymentStatusExpired).
		First(&existing).Error

	if err == nil {
		return fmt.Errorf("payment already exists for order %s", payment.OrderID)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *orderRepository) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *orderRepository) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, paymentStatusID int, paidAt *time.Time) error {
	updates := map[string]interface{}{
		"payment_status_id": paymentStatusID,
		"updated_at":        timeth.Now(),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	return r.db.WithContext(ctx).Model(&entity.Payment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *orderRepository) ListExpiredPayments(ctx context.Context) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.db.WithContext(ctx).
		Where("payment_status_id = ?", entity.PaymentStatusPending).
		Where("expires_at < ?", timeth.Now()).
		Find(&payments).Error
	return payments, err
}

func (r *orderRepository) ListDeliveredOrdersOlderThan(ctx context.Context, days int) ([]*entity.ShopOrder, error) {
	var shopOrders []*entity.ShopOrder
	cutoffTime := timeth.Now().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Preload("Order").
		Where("order_status_id = ?", entity.OrderStatusDelivered).
		Where("updated_at < ?", cutoffTime).
		Find(&shopOrders).Error

	if err != nil {
		return nil, err
	}

	return shopOrders, nil
}

func (r *orderRepository) CreateOrderLog(ctx context.Context, log *entity.OrderLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *orderRepository) GetOrderLogsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderLog, error) {
	var logs []*entity.OrderLog
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *orderRepository) GetOrderLogsByShopOrderID(ctx context.Context, shopOrderID uuid.UUID) ([]*entity.OrderLog, error) {
	var logs []*entity.OrderLog
	err := r.db.WithContext(ctx).
		Where("shop_order_id = ?", shopOrderID).
		Order("created_at ASC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
