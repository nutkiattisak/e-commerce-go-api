package delivery

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	productRepo "ecommerce-go-api/feature/product/repository"
	productUsecase "ecommerce-go-api/feature/product/usecase"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	"ecommerce-go-api/feature/shop/usecase"
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
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	resp, err := h.shopUsecase.ListShops(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
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

// ListProducts godoc
//
//	@Summary		List my products
//	@Tags			Shops
//	@Description	Get products of the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ProductListResponse
//	@Failure		401	{object}	object
//	@Failure		500	{object}	object
//	@Router			/api/shop/products [get]
func (h *ShopHandler) ListProducts(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	items, total, err := h.shopUsecase.GetProductsByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	var respItems []*entity.ProductResponse
	for _, p := range items {
		pr := &entity.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			Price:       p.Price,
			StockQty:    p.StockQty,
			IsActive:    p.IsActive,
			ShopID:      p.ShopID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
		if p.Shop.ID != uuid.Nil {
			pr.Shop = &entity.ShopResponse{ID: p.Shop.ID, Name: p.Shop.Name, ImageURL: p.Shop.ImageURL}
		}
		respItems = append(respItems, pr)
	}
	return response.Success(c, http.StatusOK, "ok", &entity.ProductListResponse{Items: respItems, Total: total})

}

// CreateProduct godoc
//
//	@Summary		Create product
//	@Tags			Shops
//	@Description	Create a new product for the authenticated user's shop
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.CreateProductRequest	true	"Create product payload"
//	@Success		201		{object}	entity.ProductResponse
//	@Failure		400		{object}	object
//	@Failure		401		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/shop/products [post]
func (h *ShopHandler) CreateProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req entity.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	product, err := h.productUsecase.CreateProduct(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusCreated, "created", product)
}

// UpdateProduct godoc
//
//	@Summary	Update product
//	@Tags		Shops
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		productId	path		int							true	"Product ID"
//	@Param		body		body		entity.CreateProductRequest	true	"Update product payload"
//	@Success	200			{object}	entity.ProductResponse
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	403			{object}	object
//	@Failure	404			{object}	object
//	@Router		/api/shop/products/{productId} [put]
func (h *ShopHandler) UpdateProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pid, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid product id")
	}

	var req entity.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	p, err := h.productUsecase.UpdateProduct(c.Request().Context(), userID, pid, &req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden")
		}
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, http.StatusNotFound, "product not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	resp := &entity.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Price:       p.Price,
		StockQty:    p.StockQty,
		IsActive:    p.IsActive,
		ShopID:      p.ShopID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Shop.ID != uuid.Nil {
		resp.Shop = &entity.ShopResponse{ID: p.Shop.ID, Name: p.Shop.Name, ImageURL: p.Shop.ImageURL}
	}

	return response.Success(c, http.StatusOK, "product updated", resp)
}

// DeleteProduct godoc
//
//	@Summary	Delete product
//	@Tags		Shops
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		productId	path		int	true	"Product ID"
//	@Success	200			{object}	object
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	403			{object}	object
//	@Failure	404			{object}	object
//	@Router		/api/shop/products/{productId} [delete]
func (h *ShopHandler) DeleteProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pid, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid product id")
	}

	if err := h.productUsecase.DeleteProduct(c.Request().Context(), userID, pid); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden")
		}
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, http.StatusNotFound, "product not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "deleted", nil)
}

// GetProduct godoc
//
//	@Summary     Get product (my shop)
//	@Tags        Shops
//	@Description Get a single product belonging to the authenticated user's shop
//	@Accept      json
//	@Produce     json
//	@Param       productId  path      int  true  "Product ID"
//	@Success     200        {object}  entity.ProductResponse
//	@Failure     400        {object}  "Bad Request"
//	@Failure     401        {object}  "Unauthorized"
//	@Failure     403        {object}  "Forbidden"
//	@Failure     404        {object}  "Not Found"
//	@Failure     500        {object}  "Internal Server Error"
//	@Router      /api/shop/products/{productId} [get]
func (h *ShopHandler) GetProduct(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pid, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid product id")
	}

	p, err := h.productUsecase.GetProductByID(c.Request().Context(), pid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, "product not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	if p.ShopID != shop.ID {
		return response.Error(c, http.StatusForbidden, "forbidden")
	}

	resp := &entity.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Price:       p.Price,
		StockQty:    p.StockQty,
		IsActive:    p.IsActive,
		ShopID:      p.ShopID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Shop.ID != uuid.Nil {
		resp.Shop = &entity.ShopResponse{ID: p.Shop.ID, Name: p.Shop.Name, ImageURL: p.Shop.ImageURL}
	}

	return response.Success(c, http.StatusOK, "ok", resp)
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
//	@Failure		400		{object}	"Bad Request"
//	@Failure		401		{object}	"Unauthorized"
//	@Failure		403		{object}	"Forbidden"
//	@Failure		500		{object}	"Internal Server Error"
//	@Router			/api/shops [put]
func (h *ShopHandler) UpdateMyShop(c echo.Context) error {
	uidKey := c.Get(middleware.CTX_KEY_USER_ID)
	uid, ok := uidKey.(uuid.UUID)
	if !ok || uid == uuid.Nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req entity.UpdateShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	shop, err := h.shopUsecase.GetShopByUserID(c.Request().Context(), uid)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	updated, err := h.shopUsecase.UpdateShop(c.Request().Context(), shop.ID, uid, &req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden")
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
