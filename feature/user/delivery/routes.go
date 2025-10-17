package delivery

import (
	"github.com/labstack/echo/v4"

	"ecommerce-go-api/middleware"
)

func (h *UserHandler) RegisterRoutes(r *echo.Group) {
	profile := r.Group("/profile", middleware.JWTAuth())
	{
		profile.GET("", h.GetProfile)
		profile.PATCH("", h.UpdateProfile)

		profile.GET("/addresses", h.GetAddresses)
		profile.POST("/addresses", h.CreateAddress)
		profile.GET("/addresses/:addressId", h.GetAddressByID)
		profile.PATCH("/addresses/:addressId", h.UpdateAddress)
		profile.DELETE("/addresses/:addressId", h.DeleteAddress)
	}
}
