package delivery

import (
	"ecommerce-go-api/middleware"

	"github.com/labstack/echo/v4"
)

func (h *ShopHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/shops", h.ListShops)
	group.GET("/shops/:shopId", h.GetShop)

	shopGroup := group.Group("")
	shopGroup.Use(middleware.JWTAuth())
	shopGroup.GET("/shop", h.GetMyShop)
	shopGroup.PUT("/shop", h.UpdateMyShop)
}
