package delivery

import (
	"ecommerce-go-api/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *OrderHandler) {
	orderGroup := g.Group("/order-groups", middleware.JWTAuth(), middleware.UserOnly())
	orderGroup.GET("", h.ListOrderGroups)
	orderGroup.GET("/:orderId", h.GetOrderGroup)

	order := g.Group("/orders", middleware.JWTAuth(), middleware.UserOnly())
	order.POST("", h.CreateOrder)
	order.GET("", h.ListOrders)
	order.GET("/:shopOrderId", h.GetOrder)
	order.POST("/:orderId/payment", h.CreateOrderPayment)
	order.GET("/:shopOrderId/tracking", h.GetShipmentTracking)

	shopOrder := g.Group("/shop/orders", middleware.JWTAuth(), middleware.ShopOwnerOnly())
	shopOrder.GET("", h.ListShopOrders)
	shopOrder.GET("/:shopOrderId", h.GetShopOrder)
	shopOrder.PUT("/:shopOrderId/status", h.UpdateShopOrderStatus)
	shopOrder.PUT("/:shopOrderId/cancel", h.CancelShopOrder)
	shopOrder.POST("/:shopOrderId/shipping", h.AddShipment)
}
