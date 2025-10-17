package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/feature/location/repository"
	"ecommerce-go-api/feature/location/usecase"
	"ecommerce-go-api/internal/response"

	"gorm.io/gorm"
)

type LocationHandler struct {
	usecase domain.LocationUsecase
}

func NewLocationHandler(u domain.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: u}
}

func RegisterLocationHandler(group *echo.Group, db *gorm.DB) {
	repo := repository.NewLocationRepository(db)
	uc := usecase.NewLocationUsecase(repo)
	h := NewLocationHandler(uc)

	h.RegisterRoutes(group)
}

// GetProvinces godoc
//
//	@Summary		Get provinces
//	@Description	Get list of provinces
//	@Tags			Location
//	@Produce		json
//	@Success		200	{array}		entity.ProvinceResponse
//	@Failure		500	{object}	object
//	@Router			/api/locations/provinces [get]
func (h *LocationHandler) GetProvinces(c echo.Context) error {
	ctx := c.Request().Context()
	provinces, err := h.usecase.GetProvinces(ctx)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	resp := make([]entity.ProvinceResponse, 0, len(provinces))
	for _, p := range provinces {
		if p == nil {
			continue
		}
		resp = append(resp, entity.ProvinceResponse{
			ID:     p.ID,
			NameTH: p.NameTH,
			NameEN: p.NameEN,
		})
	}

	return response.Success(c, http.StatusOK, "provinces retrieved", resp)
}

// GetDistrictsByProvince godoc
//
//	@Summary		Get districts by province
//	@Description	Get list of districts for a given province id
//	@Tags			Location
//	@Produce		json
//	@Param			provinceId	query		int	true	"Province ID"
//	@Success		200			{array}		entity.DistrictResponse
//	@Failure		400			{object}	object
//	@Failure		500			{object}	object
//	@Router			/api/locations/districts [get]
func (h *LocationHandler) GetDistrictsByProvince(c echo.Context) error {
	provinceId := c.QueryParam("provinceId")
	if provinceId == "" {
		return response.Error(c, http.StatusBadRequest, "provinceId is required")
	}
	id, err := strconv.Atoi(provinceId)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid provinceId")
	}

	districts, err := h.usecase.GetDistrictsByProvince(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	resp := make([]entity.DistrictResponse, 0, len(districts))
	for _, d := range districts {
		if d == nil {
			continue
		}
		var provResp *entity.ProvinceResponse
		if d.Province != nil {
			prov := d.Province
			provResp = &entity.ProvinceResponse{
				ID:     prov.ID,
				NameTH: prov.NameTH,
				NameEN: prov.NameEN,
			}
		}

		resp = append(resp, entity.DistrictResponse{
			ID:         d.ID,
			ProvinceID: d.ProvinceID,
			NameTH:     d.NameTH,
			NameEN:     d.NameEN,
			Province:   provResp,
		})
	}

	return response.Success(c, http.StatusOK, "districts retrieved", resp)
}

// GetSubDistrictsByDistrict godoc
//
//	@Summary		Get subdistricts by district
//	@Description	Get list of subdistricts for a given district id
//	@Tags			Location
//	@Produce		json
//	@Param			districtId	query		int	true	"District ID"
//	@Success		200			{array}		entity.SubDistrictResponse
//	@Failure		400			{object}	object
//	@Failure		500			{object}	object
//	@Router			/api/locations/sub-districts [get]
func (h *LocationHandler) GetSubDistrictsByDistrict(c echo.Context) error {
	districtId := c.QueryParam("districtId")
	if districtId == "" {
		return response.Error(c, http.StatusBadRequest, "districtId is required")
	}
	id, err := strconv.Atoi(districtId)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid districtId")
	}

	subs, err := h.usecase.GetSubDistrictsByDistrict(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	resp := make([]entity.SubDistrictResponse, 0, len(subs))
	for _, s := range subs {
		if s == nil {
			continue
		}

		resp = append(resp, entity.SubDistrictResponse{
			ID:         s.ID,
			Zipcode:    s.Zipcode,
			NameTH:     s.NameTH,
			NameEN:     s.NameEN,
			DistrictID: s.DistrictID,
		})
	}

	return response.Success(c, http.StatusOK, "subdistricts retrieved", resp)
}
