package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/feature/location/repository"
	"ecommerce-go-api/feature/location/usecase"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
)

type LocationHandler struct {
	usecase domain.LocationUsecase
}

func NewLocationHandler(u domain.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: u}
}

// GetProvinces godoc
//
//	@Summary		Get provinces
//	@Description	Get list of provinces
//	@Tags			Location
//	@Produce		json
//	@Success		200	{array}		entity.ProvinceResponse
//	@Failure		400	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/locations/provinces [get]
func (h *LocationHandler) GetProvinces(c echo.Context) error {
	ctx := c.Request().Context()
	provinces, err := h.usecase.GetProvinces(ctx)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "provinces retrieved", provinces)
}

// GetDistrictsByProvince godoc
//
//	@Summary		Get districts by province
//	@Description	Get list of districts for a given province id
//	@Tags			Location
//	@Produce		json
//	@Param			provinceId	query		int	true	"Province ID"
//	@Success		200			{array}		entity.DistrictResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/locations/districts [get]
func (h *LocationHandler) GetDistrictsByProvince(c echo.Context) error {
	provinceId := c.QueryParam("provinceId")
	if provinceId == "" {
		return response.Error(c, http.StatusBadRequest, errmap.ErrProvinceIDRequired.Error())
	}
	id, err := strconv.Atoi(provinceId)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidProvinceID.Error())
	}

	districts, err := h.usecase.GetDistrictsByProvince(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "districts retrieved", districts)
}

// GetSubDistrictsByDistrict godoc
//
//	@Summary		Get subdistricts by district
//	@Description	Get list of subdistricts for a given district id
//	@Tags			Location
//	@Produce		json
//	@Param			districtId	query		int	true	"District ID"
//	@Success		200			{array}		entity.SubDistrictResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/locations/sub-districts [get]
func (h *LocationHandler) GetSubDistrictsByDistrict(c echo.Context) error {
	districtId := c.QueryParam("districtId")
	if districtId == "" {
		return response.Error(c, http.StatusBadRequest, errmap.ErrDistrictIDRequired.Error())
	}
	id, err := strconv.Atoi(districtId)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidDistrictID.Error())
	}

	subs, err := h.usecase.GetSubDistrictsByDistrict(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "subdistricts retrieved", subs)
}

func RegisterLocationHandler(group *echo.Group, db *gorm.DB) {
	repo := repository.NewLocationRepository(db)
	uc := usecase.NewLocationUsecase(repo)
	h := NewLocationHandler(uc)

	h.RegisterRoutes(group)
}
