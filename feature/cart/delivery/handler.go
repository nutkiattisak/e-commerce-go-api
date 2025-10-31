package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"

	cartRepo "ecommerce-go-api/feature/cart/repository"
	cartUsecase "ecommerce-go-api/feature/cart/usecase"
	orderRepo "ecommerce-go-api/feature/order/repository"
	orderUsecase "ecommerce-go-api/feature/order/usecase"
	productRepo "ecommerce-go-api/feature/product/repository"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	userRepo "ecommerce-go-api/feature/user/repository"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type CartHandler struct {
	repo         domain.CartRepository
	cartUsecase  domain.CartUsecase
	orderUsecase domain.OrderUsecase
}

func NewCartHandler(r domain.CartRepository, u domain.CartUsecase, ou domain.OrderUsecase) *CartHandler {
	return &CartHandler{repo: r, cartUsecase: u, orderUsecase: ou}
}

// AddItem godoc
//
//	@Summary		Add item to cart
//	@Description	Adds a product to the user's cart or updates quantity if it already exists
//	@Tags			Cart
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.CartItemRequest	true	"Cart item payload"
//	@Success		201		{object}	entity.CartItemResponse	"Created"
//	@Success		200		{object}	entity.CartItemResponse	"Updated"
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		404		{object}	response.ResponseError
//	@Failure		409		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/cart [post]
func (h *CartHandler) AddItem(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.CartItemRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	item, created, err := h.cartUsecase.AddItem(c.Request().Context(), userID, req.ProductID, req.Qty)
	if err != nil {
		switch {
		case errors.Is(err, errmap.ErrQuantityMustBeGreaterThanZero),
			errors.Is(err, errmap.ErrProductInactive):
			return response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, errmap.ErrInsufficientStock):
			return response.Error(c, http.StatusConflict, err.Error())
		case errors.Is(err, errmap.ErrProductNotFound):
			return response.Error(c, http.StatusNotFound, err.Error())
		default:
			return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
		}
	}

	productSummary := entity.ProductSummary{
		ID:       item.Product.ID,
		Name:     item.Product.Name,
		ImageURL: item.Product.ImageURL,
		Price:    item.Product.Price,
		StockQty: item.Product.StockQty,
	}

	var shopResponse *entity.CartShopResponse
	if item.Product.Shop.ID != (entity.Shop{}).ID {
		shopResponse = &entity.CartShopResponse{
			ID:          item.Product.Shop.ID,
			Name:        item.Product.Shop.Name,
			Description: item.Product.Shop.Description,
			ImageURL:    item.Product.Shop.ImageURL,
		}
	}

	cartItemResponse := entity.CartItemResponse{
		ID:        item.ID,
		Product:   productSummary,
		Shop:      shopResponse,
		Qty:       item.Qty,
		UnitPrice: item.Product.Price,
		Subtotal:  float64(item.Qty) * item.Product.Price,
	}

	if created {
		return response.Success(c, http.StatusCreated, "created", cartItemResponse)
	}
	return response.Success(c, http.StatusOK, "updated", cartItemResponse)
}

// GetCart godoc
//
//	@Summary		Get user's cart
//	@Description	Retrieves the current user's cart along with items and summary
//	@Tags			Cart
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	entity.CartResponse
//	@Failure		401	{object}	response.ResponseError
//	@Failure		403	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/cart [get]
func (h *CartHandler) GetCart(c echo.Context) error {
	userID, exit := middleware.GetUserID(c)
	if exit != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}
	cart, items, summary, err := h.cartUsecase.GetCart(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	itemResponses := make([]entity.CartItemResponse, 0, len(items))
	for _, it := range items {
		unitPrice := it.Product.Price
		lineSubtotal := float64(it.Qty) * unitPrice

		pr := entity.ProductSummary{
			ID:       it.Product.ID,
			Name:     it.Product.Name,
			ImageURL: it.Product.ImageURL,
			Price:    it.Product.Price,
			StockQty: it.Product.StockQty,
		}

		var shopResponse *entity.CartShopResponse
		if it.Product.Shop.ID != (entity.Shop{}).ID {
			shopResponse = &entity.CartShopResponse{
				ID:          it.Product.Shop.ID,
				Name:        it.Product.Shop.Name,
				Description: it.Product.Shop.Description,
				ImageURL:    it.Product.Shop.ImageURL,
			}
		}

		itemResponses = append(itemResponses, entity.CartItemResponse{
			ID:        it.ID,
			Product:   pr,
			Shop:      shopResponse,
			Qty:       it.Qty,
			UnitPrice: unitPrice,
			Subtotal:  lineSubtotal,
		})
	}

	cartResponse := entity.CartResponse{
		ID:        cart.ID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		Summary:   *summary,
		Items:     itemResponses,
	}

	return response.Success(c, http.StatusOK, "ok", cartResponse)
}

// UpdateItem godoc
//
//	@Summary	Update cart item quantity
//	@Tags		Cart
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		itemId	path		int								true	"Item ID"
//	@Param		body	body		entity.UpdateCartItemRequest	true	"Update quantity payload"
//	@Success	204		{object}	object
//	@Failure	400		{object}	response.ResponseError
//	@Failure	403		{object}	response.ResponseError
//	@Failure	404		{object}	response.ResponseError
//	@Failure	409		{object}	response.ResponseError
//	@Router		/api/cart/items/{itemId} [put]
func (h *CartHandler) UpdateItem(c echo.Context) error {
	userID, exit := middleware.GetUserID(c)
	if exit != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid item id")
	}

	var req entity.UpdateCartItemRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	_, err = h.cartUsecase.UpdateItem(c.Request().Context(), userID, itemID, req.Qty)
	if err != nil {
		if errors.Is(err, errmap.ErrQuantityMustBeGreaterThanZero) {
			return response.Error(c, http.StatusBadRequest, errmap.ErrQuantityMustBeGreaterThanZero.Error())
		}
		if errors.Is(err, errmap.ErrFailedToGetCartItem) || errors.Is(err, errmap.ErrFailedToGetUserCart) {
			return response.Error(c, http.StatusNotFound, errmap.ErrFailedToGetCartItem.Error())
		}
		if errors.Is(err, errmap.ErrUnauthorized) {
			return response.Error(c, http.StatusForbidden, errmap.ErrUnauthorized.Error())
		}
		if errors.Is(err, errmap.ErrInsufficientStock) {
			return response.Error(c, http.StatusConflict, errmap.ErrInsufficientStock.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// DeleteItem godoc
//
//	@Summary	Delete cart item
//	@Tags		Cart
//	@Security	BearerAuth
//	@Param		itemId	path		int	true	"Item ID"
//	@Success	204		{object}	object
//	@Failure	400		{object}	response.ResponseError
//	@Failure	403		{object}	response.ResponseError
//	@Failure	404		{object}	response.ResponseError
//	@Router		/api/cart/items/{itemId} [delete]
func (h *CartHandler) DeleteItem(c echo.Context) error {
	userID, exit := middleware.GetUserID(c)
	if exit != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid item id")
	}

	if err := h.cartUsecase.DeleteItem(c.Request().Context(), userID, itemID); err != nil {
		if errors.Is(err, errmap.ErrFailedToGetCartItem) {
			return response.Error(c, http.StatusNotFound, errmap.ErrFailedToGetCartItem.Error())
		}
		if errors.Is(err, errmap.ErrUnauthorized) {
			return response.Error(c, http.StatusForbidden, errmap.ErrUnauthorized.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// Estimate godoc
//
//	@Summary		Estimate shipping per shop for given cart items
//	@Description	Calculate shipping costs grouped by shop for selected cart items
//	@Tags			Cart
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.EstimateShippingRequest	true	"Cart item IDs to estimate"
//	@Success		200		{object}	entity.CartShippingEstimateResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/cart/estimate [post]
func (h *CartHandler) Estimate(c echo.Context) error {
	userID, exit := middleware.GetUserID(c)
	if exit != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.EstimateShippingRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.cartUsecase.EstimateShipping(c.Request().Context(), userID, req.CartItemIDs)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, "estimate", resp)
}

func RegisterCartHandler(group *echo.Group, db *gorm.DB) {
	repo := cartRepo.NewCartRepository(db)
	shopRepository := shopRepo.NewShopRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	orderRepository := orderRepo.NewOrderRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	orderUsecase := orderUsecase.NewOrderUsecase(orderRepository, shopRepository, productRepository, userRepository)
	cartUsecase := cartUsecase.NewCartUsecase(repo, productRepository, shopRepository)
	cartHandler := NewCartHandler(repo, cartUsecase, orderUsecase)
	cartHandler.RegisterRoutes(group)
}
