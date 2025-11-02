package delivery

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type RefundHandler struct {
	usecase domain.RefundUsecase
}

func NewRefundHandler(u domain.RefundUsecase) *RefundHandler {
	return &RefundHandler{usecase: u}
}

// CreateRefund godoc
//
//	@Summary		Create refund for shop order
//	@Description	Create a refund request for a cancelled shop order
//	@Tags			Refund
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.CreateRefundRequest	true	"Refund creation payload"
//	@Success		201		{object}	entity.RefundResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		403		{object}	response.ResponseError
//	@Failure		404		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/shop/refunds [post]
func (h *RefundHandler) CreateRefund(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.CreateRefundRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	refund, err := h.usecase.CreateRefund(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusCreated, "refund created successfully", refund)
}

// ApproveRefund godoc
//
//	@Summary		Approve refund
//	@Description	Approve a refund request for shop owner
//	@Tags			Refund
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			refundId	path		string	true	"Refund ID"
//	@Success		200			{object}	entity.RefundResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		403			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/shop/refunds/{refundId}/approve [put]
func (h *RefundHandler) ApproveRefund(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	refundID, err := uuid.Parse(c.Param("refundId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	refund, err := h.usecase.ApproveRefund(c.Request().Context(), userID, refundID)
	if err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "refund approved successfully", refund)
}

// SubmitRefundBankAccount godoc
//
//	@Summary		Submit bank account for refund
//	@Description	Customer submits bank account information to receive refund
//	@Tags			Refund
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			refundId	path		string									true	"Refund ID"
//	@Param			body		body		entity.SubmitRefundBankAccountRequest	true	"Bank account information"
//	@Success		200			{object}	entity.RefundResponse
//	@Failure		400			{object}	response.ResponseError
//	@Failure		401			{object}	response.ResponseError
//	@Failure		403			{object}	response.ResponseError
//	@Failure		404			{object}	response.ResponseError
//	@Failure		500			{object}	response.ResponseError
//	@Router			/api/v1/refunds/{refundId}/bank-account [post]
func (h *RefundHandler) SubmitRefundBankAccount(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	refundID, err := uuid.Parse(c.Param("refundId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	var req entity.SubmitRefundBankAccountRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	refund, err := h.usecase.SubmitRefundBankAccount(c.Request().Context(), userID, refundID, req)
	if err != nil {
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "bank account submitted successfully", refund)
}
