package delivery

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/feature/order/repository"
	usecase "ecommerce-go-api/feature/order/usecase"
	productRepo "ecommerce-go-api/feature/product/repository"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	userRepo "ecommerce-go-api/feature/user/repository"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type OrderHandler struct {
	usecase domain.OrderUsecase
}

func NewOrderHandler(u domain.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: u}
}

// ListOrderGroups godoc
//
//	@Summary		List order groups
//	@Description	Get list of order groups (full order with all shop orders)
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			perPage		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			searchText	query		string	false	"Search by payment method or shipping name"
//	@Success		200			{object}	entity.OrderGroupListPaginationResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/order-groups [get]
func (h *OrderHandler) ListOrderGroups(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.OrderListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	orderGroups, err := h.usecase.ListOrderGroups(c.Request().Context(), userID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", orderGroups)
}

// GetOrderGroup godoc
//
//	@Summary		Get order group detail
//	@Description	Get details of a specific order group (full order with all shop orders)
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			orderGroupId	path		string	true	"Order Group ID (Main Order ID)"
//	@Success		200				{object}	entity.OrderResponse
//	@Failure		400				{object}	response.ResponseError
//	@Failure		401				{object}	response.ResponseError
//	@Failure		403				{object}	response.ResponseError
//	@Failure		404				{object}	response.ResponseError
//	@Failure		500				{object}	response.ResponseError
//	@Router			/api/order-groups/{orderGroupId} [get]
func (h *OrderHandler) GetOrderGroup(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	orderGroupID, err := uuid.Parse(c.Param("orderGroupId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	orderGroup, err := h.usecase.GetOrderGroup(c.Request().Context(), userID, orderGroupID)
	if err != nil {
		if err == errmap.ErrForbidden {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "order group not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", orderGroup)
}

// CreateOrderPayment godoc
//
//	@Summary		Create payment for order
//	@Description	Create a payment transaction
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		string						true	"Order ID"
//	@Param			body	body		entity.CreatePaymentRequest	true	"Payment payload"
//	@Success		201		{object}	entity.PaymentResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		404		{object}	response.ResponseError
//	@Failure		409		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/orders/{orderId}/payment [post]
func (h *OrderHandler) CreateOrderPayment(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	var req entity.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	payment, err := h.usecase.CreateOrderPayment(c.Request().Context(), userID, orderID, req)

	if err != nil {
		if err == errmap.ErrForbidden {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, errmap.ErrOrderNotFound.Error())
		}
		if err.Error() == "payment already exists for this order" {
			return response.Error(c, http.StatusConflict, err.Error())
		}
		if err.Error() == "payment amount does not match order total" {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "payment created successfully", payment)
}

// ListOrders godoc
//
//	@Summary		List user orders
//	@Description	Get list of orders
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			perPage		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			searchText	query		string	false	"Search by order number or status"
//	@Success		200			{object}	entity.OrderListPaginationResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/orders [get]
func (h *OrderHandler) ListOrders(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.OrderListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	orders, err := h.usecase.ListOrders(c.Request().Context(), userID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", orders)
}

// GetOrder godoc
//
//	@Summary		Get user order
//	@Description	Get details of a specific shop order
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			orderId	path		string	true	"Shop Order ID"
//	@Success		200		{object}	entity.OrderListResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/orders/{orderId} [get]
func (h *OrderHandler) GetOrder(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	shopOrderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	order, err := h.usecase.GetOrder(c.Request().Context(), userID, shopOrderID)
	if err != nil {
		if err == errmap.ErrForbidden {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", order)
}

// CreateOrder godoc
//
//	@Summary		Create a new order from cart
//	@Description	Create a new order
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.CreateOrderRequest	true	"Order creation payload"
//	@Success		201		{object}	entity.OrderResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	order, err := h.usecase.CreateOrderFromCart(c.Request().Context(), userID, req)

	if err != nil {
		switch err {
		case errmap.ErrNoShippingOptions:
			return response.Error(c, http.StatusUnprocessableEntity, errmap.ErrNoShippingOptions.Error())
		case errmap.ErrInsufficientStock:
			return response.Error(c, http.StatusConflict, errmap.ErrInsufficientStock.Error())
		default:
			return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
		}
	}

	return response.Success(c, http.StatusCreated, "created", order)
}

// ListShopOrders godoc
//
//	@Summary		List shop orders
//	@Description	Get list of orders for shop owner
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			perPage		query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			searchText	query		string	false	"Search by order number"
//	@Param			status		query		string	false	"Filter by order status"
//	@Success		200			{object}	entity.ShopOrderListPaginationResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shop/orders [get]
func (h *OrderHandler) ListShopOrders(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.OrderListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	shopOrders, err := h.usecase.ListShopOrders(c.Request().Context(), userID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", shopOrders)
}

// GetShopOrder godoc
//
//	@Summary		Get shop order
//	@Description	Get details of a specific order for shop owner
//	@Tags			Order
//	@Security		BearerAuth
//	@Produce		json
//	@Param			orderId	path		string	true	"Order ID"
//	@Success		200		{object}	entity.ShopOrderResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/orders/{orderId} [get]
func (h *OrderHandler) GetShopOrder(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	shopOrder, err := h.usecase.GetShopOrder(c.Request().Context(), userID, orderID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, "ok", shopOrder)
}

// UpdateShopOrderStatus godoc
//
//	@Summary		Update shop order status
//	@Description	Update status of a specific order for the authenticated shop owner
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		string							true	"Order ID"
//	@Param			body	body		entity.UpdateOrderStatusRequest	true	"Status update payload"
//	@Success		204		{object}	object
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/orders/{orderId}/status [put]
func (h *OrderHandler) UpdateShopOrderStatus(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	var req entity.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.usecase.UpdateShopOrderStatus(c.Request().Context(), userID, orderID, req); err != nil {
		if err.Error() == "cannot cancel order via status update, use cancel endpoint instead" {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}
		if err == errmap.ErrForbidden {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// CancelShopOrder godoc
//
//	@Summary		Cancel shop order
//	@Description	Cancel a specific order for shop owner
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		string						true	"Order ID"
//	@Param			body	body		entity.CancelOrderRequest	true	"Cancel order payload"
//	@Success		200		{object}	object
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/orders/{orderId}/cancel [put]
func (h *OrderHandler) CancelShopOrder(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	id, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidOrderID.Error())
	}

	var req entity.CancelOrderRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.usecase.CancelShopOrder(c.Request().Context(), userID, id, req); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "cancelled", nil)
}

// AddShipment godoc
//
//	@Summary		Add shipment to order
//	@Description	Add shipment details to a order for shop owner
//	@Tags			Order
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		string						true	"Order ID"
//	@Param			body	body		entity.AddShipmentRequest	true	"Shipment payload"
//	@Success		201		{object}	entity.ShipmentResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/orders/{orderId}/shipping [post]
func (h *OrderHandler) AddShipment(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid order id")
	}

	var req entity.AddShipmentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	shipment, err := h.usecase.AddShipment(c.Request().Context(), userID, orderID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "created", shipment)
}

func RegisterOrderHandler(group *echo.Group, db *gorm.DB) {
	repo := repository.NewOrderRepository(db)
	shopRepo := shopRepo.NewShopRepository(db)
	productRepo := productRepo.NewProductRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	orderUsecase := usecase.NewOrderUsecase(repo, shopRepo, productRepo, userRepo)
	handler := NewOrderHandler(orderUsecase)
	RegisterRoutes(group, handler)
}
