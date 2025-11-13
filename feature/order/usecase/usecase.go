package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/constant"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/timeth"
)

type orderUsecase struct {
	repo        domain.OrderRepository
	shopRepo    domain.ShopRepository
	productRepo domain.ProductRepository
	userRepo    domain.UserRepository
}

func NewOrderUsecase(r domain.OrderRepository, s domain.ShopRepository, p domain.ProductRepository, u domain.UserRepository) domain.OrderUsecase {
	return &orderUsecase{repo: r, shopRepo: s, productRepo: p, userRepo: u}
}

func mapToCartItemResponse(item *entity.CartItem) *entity.CartItemResponse {
	if item == nil {
		return nil
	}

	resp := &entity.CartItemResponse{
		ID:        item.ID,
		Qty:       item.Qty,
		UnitPrice: item.Product.Price,
		Subtotal:  float64(item.Qty) * item.Product.Price,
	}

	if item.Product.ID != 0 {
		resp.Product = entity.ProductSummary{
			ID:       item.Product.ID,
			Name:     item.Product.Name,
			ImageURL: item.Product.ImageURL,
			Price:    item.Product.Price,
			StockQty: item.Product.StockQty,
		}
	}

	return resp
}

func mapToShipmentResponse(s *entity.Shipment) *entity.ShipmentResponse {
	if s == nil {
		return nil
	}

	resp := &entity.ShipmentResponse{
		ID:               s.ID,
		ShopOrderID:      s.ShopOrderID,
		CourierID:        s.CourierID,
		TrackingNo:       s.TrackingNo,
		ShipmentStatusID: s.ShipmentStatusID,
		CreatedAt:        s.CreatedAt,
		ShippedAt:        s.ShippedAt,
	}

	if s.Courier.ID != 0 {
		resp.Courier = &entity.CourierListResponse{
			ID:       s.Courier.ID,
			Name:     s.Courier.Name,
			ImageURL: s.Courier.ImageURL,
			Rate:     s.Courier.Rate,
		}
	}

	if s.ShipmentStatus != nil {
		resp.ShipmentStatus = &entity.ShipmentStatusResponse{
			ID:   s.ShipmentStatus.ID,
			Code: s.ShipmentStatus.Code,
			Name: s.ShipmentStatus.Name,
		}
	}

	return resp
}

func createTimeline(logs []*entity.OrderLog) []entity.OrderTimelineItem {
	if len(logs) == 0 {
		return []entity.OrderTimelineItem{}
	}

	timeline := make([]entity.OrderTimelineItem, 0, len(logs))
	for _, log := range logs {
		timeline = append(timeline, entity.OrderTimelineItem{
			StatusID:  log.OrderStatusID,
			Note:      log.Note,
			CreatedAt: log.CreatedAt,
			CreatedBy: log.CreatedBy,
		})
	}

	return timeline
}

func (u *orderUsecase) toOrderResponseWithTimeline(ctx context.Context, order *entity.Order) *entity.OrderResponse {
	resp := &entity.OrderResponse{
		ID:                  order.ID,
		GrandTotal:          order.GrandTotal,
		ShippingName:        order.ShippingName,
		ShippingPhone:       order.ShippingPhone,
		ShippingLine1:       order.ShippingLine1,
		ShippingLine2:       order.ShippingLine2,
		ShippingSubDistrict: order.ShippingSubDistrict,
		ShippingDistrict:    order.ShippingDistrict,
		ShippingProvince:    order.ShippingProvince,
		ShippingZipcode:     order.ShippingZipcode,
		PaymentMethodID:     order.PaymentMethodID,
		ShopOrders:          make([]entity.ShopOrderResponse, 0),
	}

	logs, err := u.repo.GetOrderLogsByOrderID(ctx, order.ID)
	if err == nil {
		resp.Timeline = createTimeline(logs)
	}

	for _, so := range order.ShopOrders {
		shopOrderResp := u.toShopOrderResponseWithTimeline(ctx, &so)
		resp.ShopOrders = append(resp.ShopOrders, *shopOrderResp)
	}

	return resp
}

func (u *orderUsecase) toShopOrderResponseWithTimeline(ctx context.Context, shopOrder *entity.ShopOrder) *entity.ShopOrderResponse {
	resp := &entity.ShopOrderResponse{
		ID:            shopOrder.ID,
		OrderID:       shopOrder.OrderID,
		OrderNumber:   shopOrder.OrderNumber,
		OrderStatusID: shopOrder.OrderStatusID,
		Subtotal:      shopOrder.Subtotal,
		Shipping:      shopOrder.Shipping,
		GrandTotal:    shopOrder.GrandTotal,
		CreatedAt:     shopOrder.CreatedAt,
		UpdatedAt:     shopOrder.UpdatedAt,
		OrderItems:    make([]entity.OrderItemResponse, 0),
	}

	if shopOrder.Shop.ID != uuid.Nil {
		resp.Shop = entity.OrderShopResponse{
			ID:          shopOrder.Shop.ID,
			Name:        shopOrder.Shop.Name,
			Description: shopOrder.Shop.Description,
			ImageURL:    &shopOrder.Shop.ImageURL,
		}
	}

	logs, err := u.repo.GetOrderLogsByShopOrderID(ctx, shopOrder.ID)
	if err == nil {
		resp.Timeline = createTimeline(logs)
	}

	for _, oi := range shopOrder.OrderItems {
		itemResp := entity.OrderItemResponse{
			ID:        oi.ID,
			Qty:       oi.Qty,
			UnitPrice: oi.UnitPrice,
			Subtotal:  oi.Subtotal,
			Product: entity.OrderProductResponse{
				ID:          oi.Product.ID,
				Name:        oi.Product.Name,
				Description: oi.Product.Description,
				ImageURL:    oi.Product.ImageURL,
			},
		}
		resp.OrderItems = append(resp.OrderItems, itemResp)
	}

	return resp
}

func (u *orderUsecase) AddItemToCart(ctx context.Context, userID uuid.UUID, req entity.AddItemToCartRequest) (*entity.CartItemResponse, error) {
	if req.Qty <= 0 {
		return nil, errmap.ErrInvalidQuantity
	}

	p, err := u.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if !p.IsActive {
		return nil, errmap.ErrProductNotAvailable
	}

	cart, err := u.repo.EnsureCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	item := &entity.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Qty:       req.Qty,
	}

	updated, err := u.repo.UpsertCartItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return mapToCartItemResponse(updated), nil
}

func (u *orderUsecase) CreateOrderFromCart(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.OrderResponse, error) {
	cart, err := u.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errmap.ErrCartIsEmpty
		}
		return nil, err
	}
	items, err := u.repo.ListCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errmap.ErrCartIsEmpty
	}

	addressID := req.AddressID
	if addressID == 0 {
		return nil, errmap.ErrAddressIDRequired
	}

	order := &entity.Order{
		UserID:          userID,
		AddressID:       addressID,
		PaymentMethodID: req.PaymentMethodID,
	}

	if addr, err := u.userRepo.GetAddressByID(ctx, addressID); err == nil && addr != nil {
		order.ShippingName = addr.Name
		order.ShippingPhone = addr.PhoneNumber
		order.ShippingLine1 = addr.Line1
		order.ShippingLine2 = addr.Line2

		if addr.SubDistrict != (entity.SubDistrict{}) {
			order.ShippingSubDistrict = addr.SubDistrict.NameTH
		}
		if addr.District != (entity.District{}) {
			order.ShippingDistrict = addr.District.NameTH
		}
		if addr.Province != (entity.Province{}) {
			order.ShippingProvince = addr.Province.NameTH
		}
		order.ShippingZipcode = fmt.Sprintf("%d", addr.Zipcode)
	}

	shopItems := make(map[string][]*entity.CartItem)
	for _, ci := range items {
		shopItems[ci.Product.ShopID.String()] = append(shopItems[ci.Product.ShopID.String()], ci)
	}

	var shopOrders []*entity.ShopOrder
	orderItemsByShop := make(map[string][]*entity.OrderItem)
	var grandTotal float64

	var shopIDs []uuid.UUID
	for s := range shopItems {
		sid, _ := uuid.Parse(s)
		shopIDs = append(shopIDs, sid)
	}

	scsSlice, err := u.shopRepo.ListShopCouriersByShopIDs(ctx, shopIDs)
	if err != nil {
		return nil, err
	}

	shopCouriersMap := make(map[uuid.UUID][]*entity.ShopCourier)
	for _, sc := range scsSlice {
		shopCouriersMap[sc.ShopID] = append(shopCouriersMap[sc.ShopID], sc)
	}

	for shopIDStr, cis := range shopItems {
		so := &entity.ShopOrder{
			ShopID: uuid.Nil,
		}

		sid, _ := uuid.Parse(shopIDStr)
		so.ShopID = sid
		so.OrderNumber = fmt.Sprintf("%s-%d", constant.OrderPrefix, timeth.Now().Unix())
		so.OrderStatusID = entity.OrderStatusPending

		scs := shopCouriersMap[sid]
		if len(scs) == 0 {
			return nil, errmap.ErrNoShippingOptions
		}

		var subtotal float64
		for _, ci := range cis {
			unit := ci.Product.Price
			subtotal += float64(ci.Qty) * unit
			oi := &entity.OrderItem{
				ProductID: ci.ProductID,
				Qty:       ci.Qty,
				UnitPrice: unit,
				Subtotal:  float64(ci.Qty) * unit,
			}
			orderItemsByShop[shopIDStr] = append(orderItemsByShop[shopIDStr], oi)
		}

		so.Subtotal = subtotal
		so.Shipping = scs[0].Rate
		so.GrandTotal = so.Subtotal + so.Shipping

		grandTotal += so.GrandTotal
		shopOrders = append(shopOrders, so)
	}

	order.GrandTotal = grandTotal

	transactionID := fmt.Sprintf("TXN-%d-%s", timeth.Now().Unix(), uuid.New().String()[:8])

	expiresAt := timeth.Now().Add(24 * time.Hour)
	payment := &entity.Payment{
		TransactionID:   transactionID,
		PaymentMethodID: req.PaymentMethodID,
		PaymentStatusID: entity.PaymentStatusPending,
		Amount:          grandTotal,
		ExpiresAt:       &expiresAt,
	}

	if err := u.repo.CreateFullOrder(ctx, order, shopOrders, orderItemsByShop, payment, cart.ID, userID); err != nil {
		return nil, err
	}

	fullOrder, err := u.repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return u.toOrderResponseWithTimeline(ctx, fullOrder), nil
}

func (u *orderUsecase) ListOrders(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderListPaginationResponse, error) {
	shopOrders, total, err := u.repo.ListShopOrdersByUserID(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrderListResponse, 0, len(shopOrders))
	for _, so := range shopOrders {
		orderListResp := &entity.OrderListResponse{
			ID:                  so.ID,
			OrderID:             so.OrderID,
			OrderNumber:         so.OrderNumber,
			OrderStatusID:       so.OrderStatusID,
			Shipping:            so.Shipping,
			GrandTotal:          so.GrandTotal,
			ShippingName:        so.Order.ShippingName,
			ShippingPhone:       so.Order.ShippingPhone,
			ShippingLine1:       so.Order.ShippingLine1,
			ShippingLine2:       so.Order.ShippingLine2,
			ShippingSubDistrict: so.Order.ShippingSubDistrict,
			ShippingDistrict:    so.Order.ShippingDistrict,
			ShippingProvince:    so.Order.ShippingProvince,
			ShippingZipcode:     so.Order.ShippingZipcode,
			PaymentMethodID:     so.Order.PaymentMethodID,
			CreatedAt:           so.CreatedAt,
			UpdatedAt:           so.UpdatedAt,
			OrderItems:          make([]entity.OrderItemResponse, 0),
		}

		if so.Shop.ID != uuid.Nil {
			orderListResp.Shop = entity.OrderShopResponse{
				ID:          so.Shop.ID,
				Name:        so.Shop.Name,
				Description: so.Shop.Description,
				ImageURL:    &so.Shop.ImageURL,
			}
		}

		for _, oi := range so.OrderItems {
			itemResp := entity.OrderItemResponse{
				ID:        oi.ID,
				Qty:       oi.Qty,
				UnitPrice: oi.UnitPrice,
				Subtotal:  oi.Subtotal,
				Product: entity.OrderProductResponse{
					ID:          oi.Product.ID,
					Name:        oi.Product.Name,
					Description: oi.Product.Description,
					ImageURL:    oi.Product.ImageURL,
				},
			}
			orderListResp.OrderItems = append(orderListResp.OrderItems, itemResp)
		}

		result = append(result, orderListResp)
	}

	return &entity.OrderListPaginationResponse{
		Items: result,
		Total: total,
	}, nil
}

func (u *orderUsecase) GetOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.OrderListResponse, error) {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}

	if so.Order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	orderListResp := &entity.OrderListResponse{
		ID:                  so.ID,
		OrderID:             so.OrderID,
		OrderNumber:         so.OrderNumber,
		OrderStatusID:       so.OrderStatusID,
		Shipping:            so.Shipping,
		GrandTotal:          so.GrandTotal,
		ShippingName:        so.Order.ShippingName,
		ShippingPhone:       so.Order.ShippingPhone,
		ShippingLine1:       so.Order.ShippingLine1,
		ShippingLine2:       so.Order.ShippingLine2,
		ShippingSubDistrict: so.Order.ShippingSubDistrict,
		ShippingDistrict:    so.Order.ShippingDistrict,
		ShippingProvince:    so.Order.ShippingProvince,
		ShippingZipcode:     so.Order.ShippingZipcode,
		PaymentMethodID:     so.Order.PaymentMethodID,
		CreatedAt:           so.CreatedAt,
		UpdatedAt:           so.UpdatedAt,
		OrderItems:          make([]entity.OrderItemResponse, 0),
	}

	if so.Shop.ID != uuid.Nil {
		orderListResp.Shop = entity.OrderShopResponse{
			ID:          so.Shop.ID,
			Name:        so.Shop.Name,
			Description: so.Shop.Description,
			ImageURL:    &so.Shop.ImageURL,
		}
	}

	if so.Order.PaymentMethod != nil {
		orderListResp.PaymentMethodID = so.Order.PaymentMethod.ID
	}

	for _, oi := range so.OrderItems {
		itemResp := entity.OrderItemResponse{
			ID:        oi.ID,
			Qty:       oi.Qty,
			UnitPrice: oi.UnitPrice,
			Subtotal:  oi.Subtotal,
			Product: entity.OrderProductResponse{
				ID:          oi.Product.ID,
				Name:        oi.Product.Name,
				Description: oi.Product.Description,
				ImageURL:    oi.Product.ImageURL,
			},
		}
		orderListResp.OrderItems = append(orderListResp.OrderItems, itemResp)
	}

	logs, err := u.repo.GetOrderLogsByShopOrderID(ctx, shopOrderID)
	if err == nil {
		orderListResp.Timeline = createTimeline(logs)
	}

	return orderListResp, nil
}

func (u *orderUsecase) ListOrderGroups(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderGroupListPaginationResponse, error) {
	orders, total, err := u.repo.ListOrdersByUser(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, u.toOrderResponseWithTimeline(ctx, order))
	}

	return &entity.OrderGroupListPaginationResponse{
		Items: result,
		Total: total,
	}, nil
}

func (u *orderUsecase) GetOrderGroup(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*entity.OrderResponse, error) {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	resp := u.toOrderResponseWithTimeline(ctx, order)
	return resp, nil
}

func (u *orderUsecase) CreateOrderPayment(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req entity.CreatePaymentRequest) (*entity.PaymentResponse, error) {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	existingPayment, err := u.repo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("payment not found for this order")
	}

	if existingPayment.PaymentStatusID == entity.PaymentStatusProcessing || existingPayment.PaymentStatusID == entity.PaymentStatusCompleted {
		return nil, fmt.Errorf("payment already completed for this order")
	}

	if req.Amount != order.GrandTotal {
		return nil, fmt.Errorf("payment amount does not match order total")
	}

	now := timeth.Now()
	if err := u.repo.UpdatePaymentStatus(ctx, existingPayment.ID, entity.PaymentStatusProcessing, &now); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	updatedPayment, err := u.repo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated payment: %w", err)
	}

	response := &entity.PaymentResponse{
		ID:              updatedPayment.ID,
		OrderID:         updatedPayment.OrderID,
		TransactionID:   updatedPayment.TransactionID,
		PaymentMethodID: updatedPayment.PaymentMethodID,
		PaymentStatusID: updatedPayment.PaymentStatusID,
		Amount:          updatedPayment.Amount,
		PaidAt:          updatedPayment.PaidAt,
		ExpiresAt:       updatedPayment.ExpiresAt,
		CreatedAt:       updatedPayment.CreatedAt,
		UpdatedAt:       updatedPayment.UpdatedAt,
	}

	return response, nil
}

func (u *orderUsecase) ListShopOrders(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.ShopOrderListPaginationResponse, error) {
	shop, err := u.shopRepo.GetShopByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	shopOrders, total, err := u.repo.ListShopOrdersByShopID(ctx, shop.ID, req)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.ShopOrderListResponse, 0, len(shopOrders))
	for _, so := range shopOrders {
		orderListResp := &entity.ShopOrderListResponse{
			ID:                  so.ID,
			OrderID:             so.OrderID,
			OrderNumber:         so.OrderNumber,
			OrderStatusID:       so.OrderStatusID,
			Shipping:            so.Shipping,
			GrandTotal:          so.GrandTotal,
			ShippingName:        so.Order.ShippingName,
			ShippingPhone:       so.Order.ShippingPhone,
			ShippingLine1:       so.Order.ShippingLine1,
			ShippingLine2:       so.Order.ShippingLine2,
			ShippingSubDistrict: so.Order.ShippingSubDistrict,
			ShippingDistrict:    so.Order.ShippingDistrict,
			ShippingProvince:    so.Order.ShippingProvince,
			ShippingZipcode:     so.Order.ShippingZipcode,
			PaymentMethodID:     so.Order.PaymentMethodID,
			CreatedAt:           so.CreatedAt,
			UpdatedAt:           so.UpdatedAt,
			OrderItems:          make([]entity.OrderItemResponse, 0),
		}

		for _, oi := range so.OrderItems {
			itemResp := entity.OrderItemResponse{
				ID:        oi.ID,
				Qty:       oi.Qty,
				UnitPrice: oi.UnitPrice,
				Subtotal:  oi.Subtotal,
				Product: entity.OrderProductResponse{
					ID:          oi.Product.ID,
					Name:        oi.Product.Name,
					Description: oi.Product.Description,
					ImageURL:    oi.Product.ImageURL,
				},
			}
			orderListResp.OrderItems = append(orderListResp.OrderItems, itemResp)
		}

		result = append(result, orderListResp)
	}

	return &entity.ShopOrderListPaginationResponse{
		Items: result,
		Total: total,
	}, nil
}

func (u *orderUsecase) GetShopOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShopOrderResponse, error) {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}
	if so.ShopID == uuid.Nil {
		return nil, errmap.ErrNotFound
	}

	shop, err := u.shopRepo.GetShopByID(ctx, so.ShopID)
	if err != nil {
		return nil, err
	}
	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	resp := u.toShopOrderResponseWithTimeline(ctx, so)
	return resp, nil
}

func (u *orderUsecase) UpdateShopOrderStatus(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.UpdateOrderStatusRequest) error {
	if req.OrderStatusID == nil {
		return errmap.ErrInvalidRequest
	}

	statusID := *req.OrderStatusID

	if statusID == entity.OrderStatusCancelled {
		return fmt.Errorf("cannot cancel order")
	}

	if statusID != entity.OrderStatusProcessing &&
		statusID != entity.OrderStatusDelivered &&
		statusID != entity.OrderStatusCompleted {
		return errmap.ErrInvalidRequest
	}

	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return err
	}
	shop, err := u.shopRepo.GetShopByID(ctx, so.ShopID)
	if err != nil {
		return err
	}
	if shop.UserID != userID {
		return errmap.ErrForbidden
	}

	if err := u.repo.UpdateShopOrderStatus(ctx, shopOrderID, *req.OrderStatusID); err != nil {
		return err
	}

	if statusID == entity.OrderStatusDelivered {
		if err := u.repo.UpdateShipmentStatusByShopOrderID(ctx, shopOrderID, entity.ShipmentStatusDelivered); err != nil {
			log.Printf("[ERROR] Failed to update shipment status for shop_order_id=%s: %v", shopOrderID, err)
		}
	}

	now := timeth.Now()
	orderLog := &entity.OrderLog{
		OrderID:       so.OrderID,
		ShopOrderID:   &shopOrderID,
		OrderStatusID: statusID,
		CreatedBy:     &userID,
		CreatedAt:     &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		log.Printf("[ERROR] Failed to create order log for shop_order_id=%s, order_id=%s: %v", shopOrderID, so.OrderID, err)
	}

	return nil
}

func (u *orderUsecase) CancelShopOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.CancelOrderRequest) error {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return err
	}

	shop, err := u.shopRepo.GetShopByID(ctx, so.ShopID)
	if err != nil {
		return err
	}
	if shop.UserID != userID {
		return errmap.ErrForbidden
	}

	for _, oi := range so.OrderItems {
		if err := u.productRepo.RestoreProductStock(ctx, oi.ProductID, oi.Qty); err != nil {
			log.Printf("[ERROR] Failed to restore stock for product_id=%d, quantity=%d, shop_order_id=%s: %v", oi.ProductID, oi.Qty, shopOrderID, err)
		}
	}

	if err := u.repo.CancelShopOrder(ctx, shopOrderID, req.Reason); err != nil {
		return err
	}

	now := timeth.Now()

	orderLog := &entity.OrderLog{
		OrderID:       so.OrderID,
		ShopOrderID:   &shopOrderID,
		OrderStatusID: entity.OrderStatusCancelled,
		Note:          req.Reason,
		CreatedBy:     &userID,
		CreatedAt:     &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		log.Printf("[ERROR] Failed to create order log for cancelled order, shop_order_id=%s, order_id=%s: %v", shopOrderID, so.OrderID, err)
	}

	return nil
}

func (u *orderUsecase) AddShipment(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.AddShipmentRequest) (*entity.ShipmentResponse, error) {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}
	shop, err := u.shopRepo.GetShopByID(ctx, so.ShopID)
	if err != nil {
		return nil, err
	}
	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	existingShipment, err := u.repo.GetShipmentByShopOrderID(ctx, shopOrderID)
	if err == nil && existingShipment != nil {
		return nil, errmap.ErrShipmentAlreadyExists
	}

	now := timeth.Now()
	s := &entity.Shipment{
		ShopOrderID:      shopOrderID,
		CourierID:        req.CourierID,
		TrackingNo:       req.TrackingNo,
		ShipmentStatusID: entity.ShipmentStatusInTransit,
		ShippedAt:        &now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := u.repo.AddShipment(ctx, s); err != nil {
		return nil, err
	}

	if err := u.repo.UpdateShopOrderStatus(ctx, shopOrderID, entity.OrderStatusShipped); err != nil {
		log.Printf("[WARN] Failed to update shop order status to shipped (status=3) for shop_order_id=%s: %v", shopOrderID, err)
	}

	orderLog := &entity.OrderLog{
		OrderID:       so.OrderID,
		ShopOrderID:   &shopOrderID,
		OrderStatusID: entity.OrderStatusShipped,
		CreatedBy:     &userID,
		CreatedAt:     &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		log.Printf("[ERROR] Failed to create order log for shipment, shop_order_id=%s, order_id=%s, tracking_no=%s: %v", shopOrderID, so.OrderID, req.TrackingNo, err)
	}

	return mapToShipmentResponse(s), nil
}

func (u *orderUsecase) GetShipmentTracking(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShipmentResponse, error) {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}

	if so.Order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	shipment, err := u.repo.GetShipmentByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}

	return mapToShipmentResponse(shipment), nil
}

func (u *orderUsecase) GetShopShipmentTracking(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) (*entity.ShipmentResponse, error) {
	so, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}

	shop, err := u.shopRepo.GetShopByID(ctx, so.ShopID)
	if err != nil {
		return nil, err
	}

	if shop.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	shipment, err := u.repo.GetShipmentByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	}

	return mapToShipmentResponse(shipment), nil
}

func (u *orderUsecase) ApproveOrder(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID) error {

	shopOrder, err := u.repo.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errmap.ErrOrderNotFound
		}
		return err
	}

	order, err := u.repo.GetOrderByID(ctx, shopOrder.OrderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return errmap.ErrForbidden
	}

	if shopOrder.OrderStatusID != entity.OrderStatusDelivered {
		return fmt.Errorf("can only approve orders with DELIVERED status, current status: %d", shopOrder.OrderStatusID)
	}

	if err := u.repo.UpdateShopOrderStatus(ctx, shopOrderID, entity.OrderStatusCompleted); err != nil {
		return err
	}

	now := timeth.Now()
	orderLog := &entity.OrderLog{
		OrderID:       shopOrder.OrderID,
		ShopOrderID:   &shopOrderID,
		OrderStatusID: entity.OrderStatusCompleted,
		Note:          "Customer confirmed receipt of goods",
		CreatedBy:     &userID,
		CreatedAt:     &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		log.Printf("[ERROR] Failed to create order log for approved order, shop_order_id=%s, order_id=%s: %v", shopOrderID, shopOrder.OrderID, err)
	}

	return nil
}
