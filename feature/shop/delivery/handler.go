package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	productRepo "ecommerce-go-api/feature/product/repository"
	productUsecase "ecommerce-go-api/feature/product/usecase"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	"ecommerce-go-api/feature/shop/usecase"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type ShopHandler struct {
	shopUsecase    domain.ShopUsecase
	productUsecase domain.ProductUsecase
}

func NewShopHandler(su domain.ShopUsecase, pu domain.ProductUsecase) *ShopHandler {
	return &ShopHandler{shopUsecase: su, productUsecase: pu}
}

// GetShop godoc
//
//	@Summary		Get shop
//	@Tags			Shops
//	@Description	Get shop by ID
//	@Accept			json
//	@Produce		json
//	@Param			shopId	path		string	true	"Shop ID"
//	@Success		200		{object}	entity.Shop
//	@Router			/api/shops/{shopId} [get]
func (h *ShopHandler) GetShop(c echo.Context) error {
	id, err := uuid.Parse(c.Param("shopId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid shop id")
	}
	shop, err := h.shopUsecase.GetShopByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, "ok", shop)
}

// ListShops godoc
//
//	@Summary		List shops
//	@Tags			Shops
//	@Description	Get a paginated list of shops with optional search and sorting
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"				default(1)
//	@Param			perPage		query		int		false	"Number of items per page"	default(10)	minimum(1)	maximum(100)
//	@Param			searchText	query		string	false	"Search term for shop name"
//	@Success		200			{object}	entity.ShopListResponse
//	@Failure		400			{object}	object
//	@Failure		500			{object}	object
//	@Router			/api/shops [get]
func (h *ShopHandler) ListShops(c echo.Context) error {
	var req entity.ShopListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}
	resp, err := h.shopUsecase.ListShops(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}
	return response.Success(c, http.StatusOK, "ok", resp)
}

// GetMyShop godoc
//
//	@Summary		Get my shop
//	@Tags			Shops
//	@Description	Get the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Shop
//	@Failure		401	{object}	object
//	@Failure		500	{object}	object
//	@Router			/api/shop [get]
func (h *ShopHandler) GetMyShop(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, "ok", shop)
}

// GetShopProduct godoc
//
//	@Summary		Get product (my shop)
//	@Tags			Shops
//	@Description	Get a single product belonging to the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		int	true	"Product ID"
//	@Success		200			{object}	entity.ProductResponse
//	@Failure		400			{object}	object
//	@Failure		401			{object}	object
//	@Failure		403			{object}	object
//	@Failure		404			{object}	object
//	@Failure		500			{object}	object
//	@Router			/api/shop/products/{productId} [get]
func (h *ShopHandler) GetShopProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidProductID.Error())
	}

	product, err := h.productUsecase.GetProductByID(c.Request().Context(), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	if product.ShopID != shop.ID {
		return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
	}

	productResponse := &entity.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		Price:       product.Price,
		StockQty:    product.StockQty,
		IsActive:    product.IsActive,
		ShopID:      product.ShopID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	if product.Shop.ID != uuid.Nil {
		productResponse.Shop = &entity.ProductShopResponse{ID: product.Shop.ID, Name: product.Shop.Name, ImageURL: product.Shop.ImageURL}
	}

	return response.Success(c, http.StatusOK, "ok", productResponse)
}

// UpdateMyShop godoc
//
//	@Summary		Update my shop
//	@Tags			Shops
//	@Description	Update the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.UpdateShopRequest	true	"Update shop payload"
//	@Success		200		{object}	entity.Shop
//	@Failure		400		{object}	object
//	@Failure		401		{object}	object
//	@Failure		403		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/shops [put]
func (h *ShopHandler) UpdateMyShop(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}
	var req entity.UpdateShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	updated, err := h.shopUsecase.UpdateShop(c.Request().Context(), shop.ID, userID, &req)
	if err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, "shop updated", updated)
}

func RegisterShopHandler(group *echo.Group, db *gorm.DB) {
	shopRepository := shopRepo.NewShopRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	shopUc := usecase.NewShopUsecase(shopRepository, productRepository)
	productUc := productUsecase.NewProductUsecase(productRepository, shopRepository)
	shopHandler := NewShopHandler(shopUc, productUc)
	shopHandler.RegisterRoutes(group)
}
