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
	"ecommerce-go-api/feature/product/repository"
	"ecommerce-go-api/feature/product/usecase"
	shopRepoPkg "ecommerce-go-api/feature/shop/repository"
	"ecommerce-go-api/internal/response"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(u domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func RegisterProductHandler(group *echo.Group, db *gorm.DB) {
	repo := repository.NewProductRepository(db)
	shopRepo := shopRepoPkg.NewShopRepository(db)
	uc := usecase.NewProductUsecase(repo, shopRepo)
	h := NewProductHandler(uc)
	RegisterRoutes(group, h)
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
//	@Router			/api/products [get]
func (h *ProductHandler) ListProducts(c echo.Context) error {
	var q entity.ProductListRequest
	if err := c.Bind(&q); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}
	items, total, err := h.usecase.ListProducts(c.Request().Context(), &q)
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

// GetProduct godoc
//
//	@Summary		Get product
//	@Tags			Products
//	@Description	Get product by ID (public)
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		string	true	"Product ID"
//	@Success		200			{object}	entity.ProductResponse
//	@Router			/api/products/{productId} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid product id")
	}

	p, err := h.usecase.GetProductByID(c.Request().Context(), productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, "product not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

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

	return response.Success(c, http.StatusOK, "ok", pr)
}
