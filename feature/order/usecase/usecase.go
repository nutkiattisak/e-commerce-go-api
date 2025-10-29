package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/constant"
	"ecommerce-go-api/internal/errmap"
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

	return &entity.ShipmentResponse{
		ID:               s.ID,
		ShopOrderID:      s.ShopOrderID,
		CourierID:        s.CourierID,
		TrackingNo:       s.TrackingNo,
		ShipmentStatusID: s.ShipmentStatusID,
		CreatedAt:        s.CreatedAt,
		ShippedAt:        s.ShippedAt,
		DeliveredAt:      s.DeliveredAt,
	}
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

func toOrderResponse(order *entity.Order) *entity.OrderResponse {
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

	for _, so := range order.ShopOrders {
		resp.ShopOrders = append(resp.ShopOrders, *toShopOrderResponse(&so))
	}

	return resp
}

func toShopOrderResponse(shopOrder *entity.ShopOrder) *entity.ShopOrderResponse {
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
		so.OrderNumber = fmt.Sprintf("%s-%d", constant.OrderPrefix, time.Now().Unix())
		so.OrderStatusID = 1

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

	if err := u.repo.CreateFullOrder(ctx, order, shopOrders, orderItemsByShop, cart.ID); err != nil {
		return nil, err
	}

	now := time.Now()
	orderLog := &entity.OrderLog{
		OrderID:   order.ID,
		Note:      "Order created",
		CreatedBy: &userID,
		CreatedAt: &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		fmt.Printf("Failed to create order log: %v\n", err)
	}

	for _, so := range shopOrders {
		shopOrderLog := &entity.OrderLog{
			OrderID:     order.ID,
			ShopOrderID: &so.ID,
			Note:        "Shop order created",
			CreatedBy:   &userID,
			CreatedAt:   &now,
		}
		if err := u.repo.CreateOrderLog(ctx, shopOrderLog); err != nil {
			fmt.Printf("Failed to create shop order log: %v\n", err)
		}
	}

	fullOrder, err := u.repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(fullOrder), nil
}

func (u *orderUsecase) ListOrders(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderListPaginationResponse, error) {
	shopOrders, total, err := u.repo.ListShopOrdersByUserID(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrderListResponse, 0, len(shopOrders))
	for _, so := range shopOrders {
		orderListResp := &entity.OrderListResponse{
			ID:            so.ID,
			OrderID:       so.OrderID,
			OrderNumber:   so.OrderNumber,
			OrderStatusID: so.OrderStatusID,
			Shipping:      so.Shipping,
			GrandTotal:    so.GrandTotal,
			// CancelReason:        so.Order.CancelReason,
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

		// if so.Order.PaymentMethodID != nil {
		// 	orderListResp.PaymentMethodID = so.Order.PaymentMethodID
		// }

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
		ID:            so.ID,
		OrderID:       so.OrderID,
		OrderNumber:   so.OrderNumber,
		OrderStatusID: so.OrderStatusID,
		Shipping:      so.Shipping,
		GrandTotal:    so.GrandTotal,
		// CancelReason:        so.Order.CancelReason,
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

	return orderListResp, nil
}

func (u *orderUsecase) CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req entity.CancelOrderRequest) error {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errmap.ErrForbidden
	}

	_ = req.Reason

	return nil
}

func (u *orderUsecase) ListOrderGroups(ctx context.Context, userID uuid.UUID, req entity.OrderListRequest) (*entity.OrderGroupListPaginationResponse, error) {
	orders, total, err := u.repo.ListOrdersByUser(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, toOrderResponse(order))
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

	return toOrderResponse(order), nil
}

func (u *orderUsecase) CreateOrderPayment(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req entity.CreatePaymentRequest) (*entity.PaymentResponse, error) {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.UserID != userID {
		return nil, errmap.ErrForbidden
	}

	existingPayment, _ := u.repo.GetPaymentByOrderID(ctx, orderID)
	if existingPayment != nil && existingPayment.PaymentStatusID != 4 && existingPayment.PaymentStatusID != 5 {
		return nil, fmt.Errorf("payment already exists for this order")
	}

	if req.Amount != order.GrandTotal {
		return nil, fmt.Errorf("payment amount does not match order total")
	}

	transactionID := fmt.Sprintf("TXN-%d-%s", time.Now().Unix(), uuid.New().String()[:8])

	expiresAt := time.Now().Add(30 * time.Minute)

	paidAt := time.Now()

	payment := &entity.Payment{
		OrderID:         orderID,
		TransactionID:   transactionID,
		PaymentMethodID: req.PaymentMethodID,
		PaymentStatusID: 1,
		Amount:          req.Amount,
		PaidAt:          &paidAt,
		ExpiresAt:       &expiresAt,
	}

	if err := u.repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	response := &entity.PaymentResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		TransactionID:   payment.TransactionID,
		PaymentMethodID: payment.PaymentMethodID,
		PaymentStatusID: payment.PaymentStatusID,
		Amount:          payment.Amount,
		PaidAt:          payment.PaidAt,
		ExpiresAt:       payment.ExpiresAt,
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
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

	return toShopOrderResponse(so), nil
}

func (u *orderUsecase) UpdateShopOrderStatus(ctx context.Context, userID uuid.UUID, shopOrderID uuid.UUID, req entity.UpdateOrderStatusRequest) error {
	if req.OrderStatusID == nil || *req.OrderStatusID == 6 {
		return fmt.Errorf("cannot cancel order via status update, use cancel endpoint instead")
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

	now := time.Now()
	orderLog := &entity.OrderLog{
		OrderID:     so.OrderID,
		ShopOrderID: &shopOrderID,
		Note:        fmt.Sprintf("Status updated to %d", *req.OrderStatusID),
		CreatedBy:   &userID,
		CreatedAt:   &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		fmt.Printf("Failed to create order log: %v\n", err)
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
			fmt.Printf("Failed to restore stock for product %d: %v\n", oi.ProductID, err)
		}
	}

	if err := u.repo.CancelShopOrder(ctx, shopOrderID, req.Reason); err != nil {
		return err
	}

	now := time.Now()
	note := "Order cancelled"
	if req.Reason != "" {
		note = fmt.Sprintf("Order cancelled: %s", req.Reason)
	}
	orderLog := &entity.OrderLog{
		OrderID:     so.OrderID,
		ShopOrderID: &shopOrderID,
		Note:        note,
		CreatedBy:   &userID,
		CreatedAt:   &now,
	}
	if err := u.repo.CreateOrderLog(ctx, orderLog); err != nil {
		fmt.Printf("Failed to create order log: %v\n", err)
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

	s := &entity.Shipment{
		ShopOrderID:      shopOrderID,
		CourierID:        req.CourierID,
		TrackingNo:       req.TrackingNo,
		ShipmentStatusID: 1,
	}

	if err := u.repo.AddShipment(ctx, s); err != nil {
		return nil, err
	}

	return mapToShipmentResponse(s), nil
}
