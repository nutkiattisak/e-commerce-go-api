package delivery

import (
	"errors"
	"log"
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
//	@Summary	Add item to cart
//	@Tags		Cart
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"{\"productId\":1, \"qty\":2}"
//	@Success	201		{object}	object
//	@Router		/api/cart [post]
func (h *CartHandler) AddItem(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}
	var body struct {
		ProductID int `json:"productId"`
		Qty       int `json:"qty"`
	}
	if err := c.Bind(&body); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if body.Qty <= 0 {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidQuantity.Error())
	}

	item, created, err := h.cartUsecase.AddItem(c.Request().Context(), userID, body.ProductID, body.Qty)
	if err != nil {
		if errors.Is(err, errmap.ErrQuantityMustBeGreaterThanZero) || errors.Is(err, errmap.ErrProductInactive) {
			return response.Error(c, http.StatusBadRequest, errmap.ErrQuantityMustBeGreaterThanZero.Error())
		}
		if errors.Is(err, errmap.ErrInsufficientStock) {
			return response.Error(c, http.StatusConflict, errmap.ErrInsufficientStock.Error())
		}
		if errors.Is(err, errmap.ErrProductNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}

		log.Printf("CartHandler.AddItem error: %v", err)
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	productSummary := entity.ProductSummary{
		ID:       item.Product.ID,
		Name:     item.Product.Name,
		ImageURL: item.Product.ImageURL,
		Price:    item.Product.Price,
		StockQty: item.Product.StockQty,
	}

	cartItemResponse := entity.CartItemResponse{
		ID:        item.ID,
		Product:   productSummary,
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
//	@Summary	Get user's cart
//	@Tags		Cart
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{object}	object
//	@Router		/api/cart [get]
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

		itemResponses = append(itemResponses, entity.CartItemResponse{
			ID:        it.ID,
			Product:   pr,
			Qty:       it.Qty,
			UnitPrice: unitPrice,
			Subtotal:  lineSubtotal,
		})
	}

	cartResp := entity.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID.String(),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		Summary:   *summary,
		Items:     itemResponses,
	}

	return response.Success(c, http.StatusOK, "ok", cartResp)
}

// UpdateItem godoc
//
//	@Summary	Update cart item quantity
//	@Tags		Cart
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		itemId	path		int		true	"Item ID"
//	@Param		body	body		object	true	"{\"qty\":2}"
//	@Success	204		{object}	object
//	@Failure	400		{object}	object
//	@Failure	404		{object}	object
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
	var body struct {
		Qty int `json:"qty"`
	}
	if err := c.Bind(&body); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}
	it, err := h.cartUsecase.UpdateItem(c.Request().Context(), userID, itemID, body.Qty)
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

	productSummary := entity.ProductSummary{
		ID:       it.Product.ID,
		Name:     it.Product.Name,
		ImageURL: it.Product.ImageURL,
		Price:    it.Product.Price,
		StockQty: it.Product.StockQty,
	}

	cartItemResponse := entity.CartItemResponse{
		ID:        it.ID,
		Product:   productSummary,
		Qty:       it.Qty,
		UnitPrice: it.Product.Price,
		Subtotal:  float64(it.Qty) * it.Product.Price,
	}

	return response.Success(c, http.StatusOK, "updated", cartItemResponse)
}

// DeleteItem godoc
//
//	@Summary	Delete cart item
//	@Tags		Cart
//	@Security	ApiKeyAuth
//	@Param		itemId	path		int	true	"Item ID"
//	@Success	204		{object}	object
//	@Failure	400		{object}	object
//	@Failure	404		{object}	object
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
//	@Summary	Estimate shipping per shop for given cart items
//	@Tags		Cart
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		object	true	"{\"cartItemIds\": [1,3,5]}"
//	@Success	200		{object}	object
//	@Failure	400		{object}	object
//	@Router		/api/cart/estimate [post]
func (h *CartHandler) Estimate(c echo.Context) error {
	userID, exit := middleware.GetUserID(c)
	if exit != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var body struct {
		CartItemIDs []int `json:"cartItemIds"`
	}
	if err := c.Bind(&body); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	resp, err := h.cartUsecase.EstimateShipping(c.Request().Context(), userID, body.CartItemIDs)
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
