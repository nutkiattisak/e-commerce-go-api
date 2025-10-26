package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/internal/response"
)

type CourierHandler struct {
	usecase domain.CourierUsecase
}

func NewCourierHandler(usecase domain.CourierUsecase) *CourierHandler {
	return &CourierHandler{usecase: usecase}
}

// ListCouriers godoc
//
//	@Summary		List all couriers
//	@Description	Get a list of all available couriers (shop only)
//	@Tags			Courier
//	@Produce		json
//	@Success		200	{array}		entity.CourierListResponse
//	@Failure		401	{object}	response.ResponseError
//	@Failure		500	{object}	response.ResponseError
//	@Router			/api/couriers [get]
func (h *CourierHandler) ListCouriers(c echo.Context) error {
	ctx := c.Request().Context()

	couriers, err := h.usecase.ListCouriers(ctx)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to get couriers")
	}

	return response.Success(c, http.StatusOK, "ok", couriers)
}
