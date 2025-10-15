package delivery

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, h *ProductHandler) {
	group.GET("/products", h.ListProducts)
	group.GET("/products/:productId", h.GetProduct)
}
