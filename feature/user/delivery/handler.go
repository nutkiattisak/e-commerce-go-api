package delivery

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/feature/user/repository"
	"ecommerce-go-api/feature/user/usecase"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/response"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

// AddressResponse is a flattened response shape for addresses returned to clients
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AddressResponse struct {
	ID                int       `json:"id"`
	UserID            uuid.UUID `json:"userId"`
	Name              string    `json:"name"`
	Line1             string    `json:"line1"`
	Line2             string    `json:"line2"`
	SubDistrictID     int       `json:"subDistrictId"`
	SubDistrictNameTh string    `json:"subDistrictNameTh"`
	SubDistrictNameEn string    `json:"subDistrictNameEn"`
	DistrictNameTh    string    `json:"districtNameTh"`
	DistrictNameEn    string    `json:"districtNameEn"`
	DistrictID        int       `json:"districtId"`
	ProvinceID        int       `json:"provinceId"`
	ProvinceNameTh    string    `json:"provinceNameTh"`
	ProvinceNameEn    string    `json:"provinceNameEn"`
	Zipcode           int       `json:"zipcode"`
	PhoneNumber       string    `json:"phoneNumber"`
	IsDefault         bool      `json:"isDefault"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func mapAddressToResponse(a *entity.Address) *AddressResponse {
	if a == nil {
		return nil
	}

	var sdNameTh, sdNameEn, dNameTh, dNameEn, pNameTh, pNameEn string
	var districtID, provinceID int

	if a.SubDistrict.ID != 0 {
		sdNameTh = a.SubDistrict.NameTH
		sdNameEn = a.SubDistrict.NameEN
		if a.SubDistrict.District.ID != 0 {
			dNameTh = a.SubDistrict.District.NameTH
			dNameEn = a.SubDistrict.District.NameEN
			districtID = a.SubDistrict.District.ID
			if a.SubDistrict.District.Province.ID != 0 {
				pNameTh = a.SubDistrict.District.Province.NameTH
				pNameEn = a.SubDistrict.District.Province.NameEN
				provinceID = a.SubDistrict.District.Province.ID
			}
		}
	}

	if districtID == 0 && a.District.ID != 0 {
		dNameTh = a.District.NameTH
		dNameEn = a.District.NameEN
		districtID = a.District.ID
		if a.District.Province.ID != 0 {
			pNameTh = a.District.Province.NameTH
			pNameEn = a.District.Province.NameEN
			provinceID = a.District.Province.ID
		}
	}

	if provinceID == 0 && a.Province.ID != 0 {
		pNameTh = a.Province.NameTH
		pNameEn = a.Province.NameEN
		provinceID = a.Province.ID
	}

	return &AddressResponse{
		ID:                a.ID,
		UserID:            a.UserID,
		Name:              a.Name,
		Line1:             a.Line1,
		Line2:             a.Line2,
		SubDistrictID:     a.SubDistrictID,
		SubDistrictNameTh: sdNameTh,
		SubDistrictNameEn: sdNameEn,
		DistrictNameTh:    dNameTh,
		DistrictNameEn:    dNameEn,
		DistrictID:        districtID,
		ProvinceID:        provinceID,
		ProvinceNameTh:    pNameTh,
		ProvinceNameEn:    pNameEn,
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
//	@Router			/api/users/me [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, exist := c.Get("userId").(uuid.UUID)
	if !exist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	user, err := h.usecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "profile retrieved", user)
}

// GetUserByID returns public profile for given user id
//
//	@Summary		Get user by id
//	@Description	Get public profile by user id
//	@Tags			User
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	entity.User
//	@Failure		400		{object}	object
//	@Failure		404		{object}	object
//	@Router			/api/users/{userId} [get]
func (h *UserHandler) GetUserByID(c echo.Context) error {
	userID := c.Param("userId")
	id, err := uuid.Parse(userID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid user id")
	}

	user, err := h.usecase.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "user retrieved", user)
}

// UpdateProfile godoc
//
// @Summary     Update authenticated user's profile
// @Tags        User
// @Security    ApiKeyAuth
// @Accept      json
// @Produce     json
// @Param       body  body    entity.User  true  "User payload"
// @Success     200   {object}  entity.User
// @Failure     400   {object}  object
// @Failure     401   {object}  object
// @Failure     404   {object}  object
// @Router      /users/me [patch]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, ok := c.Get("userId").(uuid.UUID)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var user entity.User
	if err := c.Bind(&user); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	user.ID = userID

	_, err := h.usecase.UpdateProfile(c.Request().Context(), &user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "user not found")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "profile updated", nil)
}

func RegisterUserHandler(group *echo.Group, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userUsecaseInstance := usecase.NewUserUsecase(userRepository)
	userHandler := NewUserHandler(userUsecaseInstance)

	userHandler.RegisterRoutes(group)
}

// GetAddresses godoc
//
//	@Summary	Get authenticated user's addresses
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{array}		AddressResponse
//	@Failure	401	{object}	object
//	@Router		/users/me/addresses [get]
func (h *UserHandler) GetAddresses(c echo.Context) error {
	userID, exist := c.Get("userId").(uuid.UUID)
	if !exist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addrs, err := h.usecase.GetAddresses(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	var out []*AddressResponse
	for _, a := range addrs {
		out = append(out, mapAddressToResponse(a))
	}

	return response.Success(c, http.StatusOK, "addresses retrieved", out)
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
//	@Router		/users/me/addresses [post]
func (h *UserHandler) CreateAddress(c echo.Context) error {
	userID, exist := c.Get("userId").(uuid.UUID)
	if !exist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var addr entity.Address
	if err := c.Bind(&addr); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	addr.UserID = userID

	_, err := h.usecase.CreateAddress(c.Request().Context(), &addr)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "address created", nil)
}

// GetAddressByID godoc
//
//	@Summary	Get address by id for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	200			{object}	AddressResponse
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	404			{object}	object
//	@Router		/users/me/addresses/{addressId} [get]
func (h *UserHandler) GetAddressByID(c echo.Context) error {
	_, exist := c.Get("userId").(uuid.UUID)
	if !exist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addressID, err := strconv.Atoi(c.Param("addressId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid address id")
	}

	addr, err := h.usecase.GetAddressByID(c.Request().Context(), addressID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "address retrieved", mapAddressToResponse(addr))
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
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	404			{object}	object
//	@Router		/users/me/addresses/{addressId} [patch]
func (h *UserHandler) UpdateAddress(c echo.Context) error {
	userID, exist := c.Get("userId").(uuid.UUID)
	if !exist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	addressID, err := strconv.Atoi(c.Param("addressId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid address id")
	}

	var addr entity.Address
	if err := c.Bind(&addr); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	addr.ID = addressID

	_, err = h.usecase.UpdateAddress(c.Request().Context(), &addr, userID)

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "address updated", nil)
}

// DeleteAddress godoc
//
//	@Summary	Soft delete address for authenticated user
//	@Tags		User
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		addressId	path		int	true	"Address ID"
//	@Success	204			{object}	object
//	@Failure	400			{object}	object
//	@Failure	401			{object}	object
//	@Failure	404			{object}	object
//	@Router		/users/me/addresses/{addressId} [delete]
func (h *UserHandler) DeleteAddress(c echo.Context) error {
	userID, exist := c.Get("userId").(uuid.UUID)
	if !exist {
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

	return response.Success(c, http.StatusNoContent, "address deleted", nil)
}
