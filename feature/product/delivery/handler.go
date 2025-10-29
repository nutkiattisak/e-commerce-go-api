package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/feature/product/repository"
	"ecommerce-go-api/feature/product/usecase"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	shopUsecase "ecommerce-go-api/feature/shop/usecase"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type ProductHandler struct {
	usecase     domain.ProductUsecase
	shopUsecase domain.ShopUsecase
}

func NewProductHandler(u domain.ProductUsecase, su domain.ShopUsecase) *ProductHandler {
	return &ProductHandler{usecase: u, shopUsecase: su}
}

// ListProducts godoc
//
//	@Summary		List products
//	@Tags			Products
//	@Description	Get public product listing with filters and pagination
//	@Accept			json
//	@Produce		json
//	@Param			searchText	query		string	false	"searchText query"
//	@Param			page		query		int		false	"page"
//	@Param			perPage		query		int		false	"perPage"
//	@Success		200			{object}	entity.ProductListResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/products [get]
func (h *ProductHandler) ListProducts(c echo.Context) error {
	var q entity.ProductListRequest
	if err := c.Bind(&q); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	result, err := h.usecase.ListProducts(c.Request().Context(), &q)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", result)
}

// GetProduct godoc
//
//	@Summary		Get product
//	@Tags			Products
//	@Description	Get product by ID (public)
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		string	true	"Product ID"
//	@Success		200			{object}	entity.ProductResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/products/{productId} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidProductID.Error())
	}

	product, err := h.usecase.GetProductByID(c.Request().Context(), productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", product)
}

// ListShopProducts godoc
//
//	@Summary		List shop products
//	@Tags			Products
//	@Description	Get products for the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	entity.ProductListResponse
//	@Failure		400	{object}	response.ResponseError
//	@Failure		401	{object}	response.ResponseError
//	@Failure		403	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/shop/products [get]
func (h *ProductHandler) ListShopProducts(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	result, err := h.usecase.GetProductsByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "ok", result)
}

// CreateShopProduct godoc
//
//	@Summary		Create product
//	@Tags			Products
//	@Description	Create a new product for the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		entity.CreateProductRequest	true	"Create Product Request"
//	@Success		201		{object}	entity.ProductResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/products [post]
func (h *ProductHandler) CreateShopProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}
	product, err := h.usecase.CreateProduct(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusCreated, "created", product)
}

// GetShopProduct godoc
//
//	@Summary		Get product (my shop)
//	@Description	Get a single product belonging to the authenticated user's shop
//	@Tags			Shops
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		int	true	"Product ID"
//	@Success		200			{object}	entity.ProductResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		403			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shop/products/{productId} [get]
func (h *ProductHandler) GetShopProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidProductID.Error())
	}

	product, err := h.usecase.GetProductByID(c.Request().Context(), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}
	if product.ShopID != shop.ID {
		return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
	}

	return response.Success(c, http.StatusOK, "ok", product)
}

// UpdateShopProduct godoc
//
//	@Summary		Update product (my shop)
//	@Tags			Shops
//	@Security		BearerAuth
//	@Description	Update a product belonging to the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		int							true	"Product ID"
//	@Param			body		body		entity.UpdateProductRequest	true	"Update Product Request"
//	@Success		200			{object}	entity.ProductResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		403			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shop/products/{productId} [put]
func (h *ProductHandler) UpdateShopProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidProductID.Error())
	}

	var req entity.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	product, err := h.usecase.UpdateProduct(c.Request().Context(), userID, productID, &req)
	if err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		if errors.Is(err, errmap.ErrProductNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}
		c.Logger().Error("Unhandled error in UpdateShopProduct: ", err)
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, "updated", product)
}

// DeleteShopProduct godoc
//
//	@Summary		Delete product (my shop)
//	@Tags			Shops
//	@Security		BearerAuth
//	@Description	Delete a product belonging to the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		int	true	"Product ID"
//	@Success		204			{object}	object
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		403			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shop/products/{productId} [delete]
func (h *ProductHandler) DeleteShopProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid product id")
	}

	if err := h.usecase.DeleteProduct(c.Request().Context(), userID, productID); err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		if errors.Is(err, errmap.ErrProductNotFound) {
			return response.Error(c, http.StatusNotFound, errmap.ErrProductNotFound.Error())
		}
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}

func RegisterProductHandler(group *echo.Group, db *gorm.DB) {
	productRepo := repository.NewProductRepository(db)
	shopRepository := shopRepo.NewShopRepository(db)

	productUsecase := usecase.NewProductUsecase(productRepo, shopRepository)
	shopUsecase := shopUsecase.NewShopUsecase(shopRepository, productRepo)

	handler := NewProductHandler(productUsecase, shopUsecase)
	RegisterRoutes(group, handler)
}
