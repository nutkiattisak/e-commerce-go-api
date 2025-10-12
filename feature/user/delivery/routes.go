package delivery

import (
	"github.com/labstack/echo/v4"

	"ecommerce-go-api/middleware"
)

func (h *UserHandler) RegisterRoutes(r *echo.Group) {
	users := r.Group("/users")
	{
		users.GET("/me", h.GetProfile, middleware.JWTAuth())
		users.GET("/:userId", h.GetUserByID, middleware.JWTAuth())
	}
}
