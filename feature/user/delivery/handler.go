package delivery

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/feature/user/repository"
	"ecommerce-go-api/feature/user/usecase"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
	"ecommerce-go-api/middleware"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

// GetProfile returns profile of authenticated user
//
//	@Summary		Get current user's profile
//	@Description	Get profile of the authenticated user
//	@Tags			User
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	entity.UserResponse
//	@Failure		400	{object}	response.ResponseError
//	@Failure		401	{object}	response.ResponseError
//	@Router			/api/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	user, err := h.usecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "profile retrieved", user)
}

// UpdateProfile godoc
//
//	@Summary	Update authenticated user's profile
//	@Tags		User
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		entity.UpdateProfileRequest	true	"Update profile payload"
//	@Success	204		{object}	object
//	@Failure	400		{object}	response.ResponseError
//	@Failure	401		{object}	response.ResponseError
//	@Failure	404		{object}	response.ResponseError
//	@Router		/api/profile [patch]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	var req entity.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	user := &entity.User{
		ID:          userID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		ImageURL:    req.ImageURL,
		UpdatedAt:   time.Now(),
	}

	if err := h.usecase.UpdateProfile(c.Request().Context(), user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "user not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// GetAddresses godoc
//
//	@Summary	Get authenticated user's addresses
//	@Tags		User
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}		entity.AddressResponse
//	@Failure	401	{object}	object
//	@Router		/api/profile/addresses [get]
func (h *UserHandler) GetAddresses(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	addrs, err := h.usecase.GetAddresses(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "addresses retrieved", addrs)
}

// GetAddressByID godoc
//
//	@Summary	Get address by id for authenticated user
//	@Tags		User
//	@Security	BearerAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	200			{object}	entity.AddressResponse
//	@Failure	400			{object}	response.ResponseError
//	@Failure	401			{object}	response.ResponseError
//	@Failure	404			{object}	response.ResponseError
//	@Router		/api/profile/addresses/{addressId} [get]
func (h *UserHandler) GetAddressByID(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}

	addressID, err := strconv.Atoi(c.Param("addressId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid address id")
	}

	addr, err := h.usecase.GetAddressByID(c.Request().Context(), addressID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "address not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	if addr == nil {
		return response.Error(c, http.StatusNotFound, "address not found")
	}

	if addr.UserID != userID {
		return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
	}

	return response.Success(c, http.StatusOK, "address retrieved", addr)
}

// CreateAddress godoc
//
//	@Summary	Create address for authenticated user
//	@Tags		User
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		entity.CreateAddressRequest	true	"Address payload"
//	@Success	201		{object}	entity.AddressResponse
//	@Failure	400		{object}	response.ResponseError
//	@Failure	401		{object}	response.ResponseError
//	@Router		/api/profile/addresses [post]
func (h *UserHandler) CreateAddress(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req entity.CreateAddressRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	addr := &entity.Address{
		UserID:        userID,
		Name:          req.Name,
		Line1:         req.Line1,
		Line2:         req.Line2,
		SubDistrictID: req.SubDistrictID,
		DistrictID:    req.DistrictID,
		ProvinceID:    req.ProvinceID,
		Zipcode:       req.Zipcode,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     req.IsDefault,
	}

	added, err := h.usecase.CreateAddress(c.Request().Context(), addr)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "address created", added)
}

// UpdateAddress godoc
//
//	@Summary	Update address for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		addressId	path		int							true	"Address ID"
//	@Param		body		body		entity.UpdateAddressRequest	true	"Address payload"
//	@Success	200			{object}	entity.AddressResponse
//	@Failure	400			{object}	response.ResponseError
//	@Failure	401			{object}	response.ResponseError
//	@Failure	404			{object}	response.ResponseError
//	@Router		/api/profile/addresses/{addressId} [patch]
func (h *UserHandler) UpdateAddress(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addressID, err := strconv.Atoi(c.Param("addressId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid address id")
	}

	var req entity.UpdateAddressRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	addr, err := h.usecase.GetAddressByID(c.Request().Context(), addressID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	if addr == nil {
		return response.Error(c, http.StatusNotFound, "address not found")
	}
	if addr.UserID != userID {
		return response.Error(c, http.StatusForbidden, "forbidden")
	}

	addrEntity := &entity.Address{
		ID:            addressID,
		UserID:        userID,
		Name:          req.Name,
		Line1:         req.Line1,
		Line2:         req.Line2,
		SubDistrictID: req.SubDistrictID,
		DistrictID:    req.DistrictID,
		ProvinceID:    req.ProvinceID,
		Zipcode:       req.Zipcode,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     req.IsDefault,
	}

	_, err = h.usecase.UpdateAddress(c.Request().Context(), addrEntity, userID)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// DeleteAddress godoc
//
//	@Summary	Soft delete address for authenticated user
//	@Tags		User
//	@Security	BearerAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	204			{object}	object
//	@Failure	400			{object}	response.ResponseError
//	@Failure	404			{object}	response.ResponseError
//	@Failure	500			{object}	response.ResponseError
//	@Router		/api/profile/addresses/{addressId} [delete]
func (h *UserHandler) DeleteAddress(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addressID, err := strconv.Atoi(c.Param("addressId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid address id")
	}

	if err := h.usecase.DeleteAddress(c.Request().Context(), addressID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusNotFound, "address not found")
		}
		if errors.Is(err, errmap.ErrForbidden) {
			return response.Error(c, http.StatusForbidden, "forbidden")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

func RegisterUserHandler(group *echo.Group, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userUsecaseInstance := usecase.NewUserUsecase(userRepository)
	userHandler := NewUserHandler(userUsecaseInstance)
	userHandler.RegisterRoutes(group)
}
