package delivery

import (
	"ecommerce-go-api/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *OrderHandler) {
	orderGroup := g.Group("/order-groups", middleware.JWTAuth(), middleware.UserOnly())
	orderGroup.GET("", h.ListOrderGroups)
	orderGroup.GET("/:orderGroupId", h.GetOrderGroup)

	order := g.Group("/orders", middleware.JWTAuth(), middleware.UserOnly())
	order.POST("", h.CreateOrder)
	order.GET("", h.ListOrders)
	order.GET("/:orderId", h.GetOrder)
	order.POST("/:orderId/payment", h.CreateOrderPayment)

	shopOrder := g.Group("/shop/orders", middleware.JWTAuth(), middleware.ShopOwnerOnly())
	shopOrder.GET("", h.ListShopOrders)
	shopOrder.GET("/:orderId", h.GetShopOrder)
	shopOrder.PUT("/:orderId/status", h.UpdateShopOrderStatus)
	shopOrder.PUT("/:orderId/cancel", h.CancelShopOrder)
	shopOrder.POST("/:orderId/shipping", h.AddShipment)
}
