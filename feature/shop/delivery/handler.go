package delivery

import (
	"errors"
	"net/http"

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
//	@Success		200		{object}	entity.ShopResponse
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
//	@Failure		400			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shops [get]
func (h *ShopHandler) ListShops(c echo.Context) error {
	var req entity.ShopListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.shopUsecase.ListShops(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, errmap.ErrNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrNotFound.Error())
		}

		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}
	return response.Success(c, http.StatusOK, "ok", resp)
}

// GetMyShop godoc
//
//	@Summary		Get my shop
//	@Description	Get the authenticated user's shop
//	@Tags			Shops
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ShopResponse
//	@Failure		401	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/shop [get]
func (h *ShopHandler) GetMyShop(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, errmap.ErrNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrNotFound.Error())
		}
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}

		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}
	return response.Success(c, http.StatusOK, "ok", shop)
}

// UpdateMyShop godoc
//
//	@Summary		Update my shop
//	@Description	Update the authenticated user's shop
//	@Tags			Shops
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.UpdateShopRequest	true	"Update shop payload"
//	@Success		200		{object}	entity.ShopResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shops [put]
func (h *ShopHandler) UpdateMyShop(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}
	var req entity.UpdateShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
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

// UpdateShopCouriers godoc
//
//	@Summary		Update shop courier
//	@Description	Update courier settings for the authenticated user's shop (soft deletes old record and creates new one)
//	@Tags			Shops
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.UpdateShopCouriersRequest	true	"Courier payload"
//	@Success		200		{object}	entity.ShopCourierResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/couriers [put]
func (h *ShopHandler) UpdateShopCouriers(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.UpdateShopCouriersRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.shopUsecase.UpdateShopCouriers(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "courier updated", resp)
}

// GetShopCouriers godoc
//
//	@Summary		Get shop courier
//	@Description	Get active courier settings for the authenticated user's shop (deleted_at IS NULL)
//	@Tags			Shops
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	entity.ShopCourierResponse
//	@Failure		401	{object}	response.ResponseError
//	@Failure		404	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/shop/couriers [get]
func (h *ShopHandler) GetShopCouriers(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	resp, err := h.shopUsecase.GetShopCouriers(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", resp)
}

func RegisterShopHandler(group *echo.Group, db *gorm.DB) {
	shopRepository := shopRepo.NewShopRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	shopUc := usecase.NewShopUsecase(shopRepository, productRepository)
	productUc := productUsecase.NewProductUsecase(productRepository, shopRepository)
	shopHandler := NewShopHandler(shopUc, productUc)
	shopHandler.RegisterRoutes(group)
}
