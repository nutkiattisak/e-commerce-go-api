package delivery

import (
	"ecommerce-go-api/middleware"

	"github.com/labstack/echo/v4"
)

func (h *CartHandler) RegisterRoutes(r *echo.Group) {
	cart := r.Group("/cart")
	cart.GET("", h.GetCart, middleware.JWTAuth(), middleware.UserOnly())
	cart.POST("", h.AddItem, middleware.JWTAuth(), middleware.UserOnly())
	cart.POST("/estimate", h.Estimate, middleware.JWTAuth(), middleware.UserOnly())
	cart.PUT("/items/:itemId", h.UpdateItem, middleware.JWTAuth(), middleware.UserOnly())
	cart.DELETE("/items/:itemId", h.DeleteItem, middleware.JWTAuth(), middleware.UserOnly())
}
