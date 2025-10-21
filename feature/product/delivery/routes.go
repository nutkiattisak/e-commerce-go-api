package delivery

import (
	"ecommerce-go-api/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, h *ProductHandler) {
	group.GET("/products", h.ListProducts)
	group.GET("/products/:productId", h.GetProduct)

	shopGroup := group.Group("/shop")
	shopGroup.Use(middleware.JWTAuth())
	shopGroup.GET("/products", h.ListShopProducts)
	shopGroup.POST("/products", h.CreateShopProduct)
	shopGroup.GET("/products/:productId", h.GetShopProduct)
	shopGroup.PUT("/products/:productId", h.UpdateShopProduct)
	shopGroup.DELETE("/products/:productId", h.DeleteShopProduct)
}
