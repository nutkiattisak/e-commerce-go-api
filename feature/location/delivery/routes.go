package delivery

import (
	"github.com/labstack/echo/v4"
)

func (h *LocationHandler) RegisterRoutes(r *echo.Group) {
	locations := r.Group("/locations")
	locations.GET("/provinces", h.GetProvinces)
	locations.GET("/districts", h.GetDistrictsByProvince)
	locations.GET("/sub-districts", h.GetSubDistrictsByDistrict)
}
