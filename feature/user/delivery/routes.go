package delivery

import (
	"github.com/labstack/echo/v4"

	"ecommerce-go-api/middleware"
)

func (h *UserHandler) RegisterRoutes(r *echo.Group) {
	users := r.Group("/users")
	{
		users.GET("/me", h.GetProfile, middleware.JWTAuth())
		users.PATCH("/me", h.UpdateProfile, middleware.JWTAuth())
		users.GET(":userId", h.GetUserByID, middleware.JWTAuth())
		users.GET("/me/addresses", h.GetAddresses, middleware.JWTAuth())
		users.POST("/me/addresses", h.CreateAddress, middleware.JWTAuth())
		users.GET("/me/addresses/:addressId", h.GetAddressByID, middleware.JWTAuth())
		users.PATCH("/me/addresses/:addressId", h.UpdateAddress, middleware.JWTAuth())
		users.DELETE("/me/addresses/:addressId", h.DeleteAddress, middleware.JWTAuth())
	}
}
