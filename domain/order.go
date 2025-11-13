package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"ecommerce-go-api/entity"
)

type OrderUsecase interface {
	CreateOrderFromCart(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.OrderResponse, error)
	AddItemToCart(ctx context.Context, userID uuid.UUID, req entity.AddItemToCartRequest) (*entity.CartItemResponse, error)
	ListOrders(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderListPaginationResponse, error)
	GetOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.OrderListResponse, error)

	ListOrderGroups(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderGroupListPaginationResponse, error)
	GetOrderGroup(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*entity.OrderResponse, error)

	CreateOrderPayment(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req entity.CreatePaymentRequest) (*entity.PaymentResponse, error)

	ListShopOrders(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.ShopOrderListPaginationResponse, error)
	GetShopOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShopOrderResponse, error)
	UpdateShopOrderStatus(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.UpdateOrderStatusRequest) error
	CancelShopOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.CancelOrderRequest) error
	AddShipment(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.AddShipmentRequest) (*entity.ShipmentResponse, error)
	GetShipmentTracking(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShipmentResponse, error)
	GetShopShipmentTracking(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShipmentResponse, error)
	ApproveOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) error
}

type OrderRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	ListCartItems(ctx context.Context, cartID int) ([]*entity.CartItem, error)
	ClearCart(ctx context.Context, cartID int) error
	EnsureCartForUser(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	AddCartItem(ctx context.Context, item *entity.CartItem) error
	UpsertCartItem(ctx context.Context, item *entity.CartItem) (*entity.CartItem, error)
	GetCartItemByID(ctx context.Context, id int) (*entity.CartItem, error)

	CreateOrder(ctx context.Context, order *entity.Order) error
	CreateShopOrder(ctx context.Context, so *entity.ShopOrder) error
	CreateOrderItems(ctx context.Context, items []*entity.OrderItem) error
	CreateFullOrder(ctx context.Context, order *entity.Order, shopOrders []*entity.ShopOrder, orderItemsByShop map[string][]*entity.OrderItem, payment *entity.Payment, cartID int, userID uuid.UUID) error

	ListOrdersByUser(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) ([]*entity.Order, int64, error)
	ListShopOrdersByUserID(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) ([]*entity.ShopOrder, int64, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)

	ListShopOrdersByShopID(ctx context.Context, shopID uuid.UUID, req entity.OrderListRequest) ([]*entity.ShopOrder, int64, error)
	GetShopOrderByID(ctx context.Context, id uuid.UUID) (*entity.ShopOrder, error)
	UpdateShopOrderStatus(ctx context.Context, id uuid.UUID, OrderStatusID int) error
	CancelShopOrder(ctx context.Context, id uuid.UUID, reason string) error

	AddShipment(ctx context.Context, s *entity.Shipment) error
	GetShipmentByShopOrderID(ctx context.Context, shopOrderID uuid.UUID) (*entity.Shipment, error)
	UpdateShipmentStatusByShopOrderID(ctx context.Context, shopOrderID uuid.UUID, shipmentStatusID int) error

	// Payment
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)
	GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, paymentStatusID int, paidAt *time.Time) error
	ListExpiredPayments(ctx context.Context) ([]*entity.Payment, error)

	// OrderLog
	CreateOrderLog(ctx context.Context, log *entity.OrderLog) error
	GetOrderLogsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderLog, error)
	GetOrderLogsByShopOrderID(ctx context.Context, shopOrderID uuid.UUID) ([]*entity.OrderLog, error)

	ListDeliveredOrdersOlderThan(ctx context.Context, days int) ([]*entity.ShopOrder, error)
}
