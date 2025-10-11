package delivery

import (
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) RegisterRoutes(r *echo.Group) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/register/shop", h.RegisterShop)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
	}
}
