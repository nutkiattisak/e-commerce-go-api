package delivery

import (
	"errors"
	"net/http"
	"strconv"

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

func mapAddressToResponse(a *entity.Address) *entity.AddressResponse {
	if a == nil {
		return nil
	}
	subDistrictNameTh, subDistrictNameEn := "", ""
	districtNameTh, districtNameEn := "", ""
	provinceNameTh, provinceNameEn := "", ""

	if a.SubDistrict != (entity.SubDistrict{}) {
		subDistrictNameTh = a.SubDistrict.NameTH
		subDistrictNameEn = a.SubDistrict.NameEN
		if a.SubDistrict.District != nil {
			districtNameTh = a.SubDistrict.District.NameTH
			districtNameEn = a.SubDistrict.District.NameEN
			if a.SubDistrict.District.Province != nil {
				provinceNameTh = a.SubDistrict.District.Province.NameTH
				provinceNameEn = a.SubDistrict.District.Province.NameEN
			}
		}
	}

	if districtNameTh == "" && a.District != (entity.District{}) {
		districtNameTh = a.District.NameTH
		districtNameEn = a.District.NameEN
	}
	if provinceNameTh == "" && a.Province != (entity.Province{}) {
		provinceNameTh = a.Province.NameTH
		provinceNameEn = a.Province.NameEN
	}

	return &entity.AddressResponse{
		ID:                a.ID,
		UserID:            a.UserID,
		Name:              a.Name,
		Line1:             a.Line1,
		Line2:             a.Line2,
		SubDistrictID:     a.SubDistrictID,
		SubDistrictNameTh: subDistrictNameTh,
		SubDistrictNameEn: subDistrictNameEn,
		DistrictNameTh:    districtNameTh,
		DistrictNameEn:    districtNameEn,
		DistrictID:        a.DistrictID,
		ProvinceID:        a.ProvinceID,
		ProvinceNameTh:    provinceNameTh,
		ProvinceNameEn:    provinceNameEn,
		Zipcode:           a.Zipcode,
		PhoneNumber:       a.PhoneNumber,
		IsDefault:         a.IsDefault,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

// GetProfile returns profile of authenticated user
//
//	@Summary		Get current user's profile
//	@Description	Get profile of the authenticated user
//	@Tags			User
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Failure		400	{object}	object
//	@Failure		401	{object}	object
//	@Router			/api/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	user, err := h.usecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Success(c, http.StatusOK, "profile retrieved", nil)
	}
	return response.Success(c, http.StatusOK, "profile retrieved", &entity.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ImageURL:    user.ImageURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

// UpdateProfile godoc
//
//	@Summary	Update authenticated user's profile
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		entity.User	true	"User payload"
//	@Success	200		{object}	entity.User
//	@Failure	400		{object}	object
//	@Failure	401		{object}	object
//	@Failure	404		{object}	object
//	@Router		/api/profile [patch]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req entity.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	user := &entity.User{
		ID:          userID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		ImageURL:    req.ImageURL,
	}

	_, exit := h.usecase.UpdateProfile(c.Request().Context(), user)
	if exit != nil {
		if exit == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "user not found")
		}
		return response.Error(c, http.StatusInternalServerError, exit.Error())
	}

	return response.NoContent(c)
}

// GetAddresses godoc
//
//	@Summary	Get authenticated user's addresses
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{array}		entity.AddressResponse
//	@Failure	401	{object}	object
//	@Router		/api/profile/addresses [get]
func (h *UserHandler) GetAddresses(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addrs, err := h.usecase.GetAddresses(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	var out []*entity.AddressResponse
	for _, a := range addrs {
		out = append(out, mapAddressToResponse(a))
	}

	return response.Success(c, http.StatusOK, "addresses retrieved", out)
}

// GetAddressByID godoc
//
//	@Summary	Get address by id for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	200			{object}	entity.AddressResponse
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	404			{object}	object
//	@Router		/api/profile/addresses/{addressId} [get]
func (h *UserHandler) GetAddressByID(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
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
		return response.Error(c, http.StatusForbidden, "forbidden")
	}

	return response.Success(c, http.StatusOK, "address retrieved", mapAddressToResponse(addr))
}

// CreateAddress godoc
//
//	@Summary	Create address for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		entity.Address	true	"Address payload"
//	@Success	201		{object}	object
//	@Failure	400		{object}	object
//	@Failure	401		{object}	object
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

	full, err := h.usecase.GetAddressByID(c.Request().Context(), added.ID)
	if err != nil {
		return response.Success(c, http.StatusCreated, "address created", mapAddressToResponse(added))
	}

	return response.Success(c, http.StatusCreated, "address created", mapAddressToResponse(full))
}

// UpdateAddress godoc
//
//	@Summary	Update address for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Accept		json
//	@Produce	json
//	@Param		addressId	path		int				true	"Address ID"
//	@Param		body		body		entity.Address	true	"Address payload"
//	@Success	200			{object}	entity.Address
//	@Failure	400			{object}	"Bad Request"
//	@Failure	401			{object}	"Unauthorized"
//	@Failure	404			{object}	"Not Found"
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

	addr.Name = req.Name
	addr.Line1 = req.Line1
	addr.Line2 = req.Line2
	addr.SubDistrictID = req.SubDistrictID
	addr.DistrictID = req.DistrictID
	addr.ProvinceID = req.ProvinceID
	addr.Zipcode = req.Zipcode
	addr.PhoneNumber = req.PhoneNumber
	addr.IsDefault = req.IsDefault

	_, err = h.usecase.UpdateAddress(c.Request().Context(), addr, userID)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c)
}

// DeleteAddress godoc
//
//	@Summary	Soft delete address for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	204			{object}	object
//	@Failure		400			{object}	"Bad Request"
//	@Failure		404			{object}	"Not Found"
//	@Failure		500			{object}	"Internal Server Error"
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
