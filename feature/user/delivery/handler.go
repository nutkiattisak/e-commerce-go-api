package delivery

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/feature/user/repository"
	"ecommerce-go-api/feature/user/usecase"
	"ecommerce-go-api/internal/response"
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
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Failure		400	{object}	object
//	@Failure		401	{object}	object
//	@Router			/api/users/me [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	uid := c.Get("user_id")
	if uid == nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var userID uuid.UUID
	switch v := uid.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid user id")
		}
		userID = parsed
	case uuid.UUID:
		userID = v
	default:
		return response.Error(c, http.StatusBadRequest, "invalid user id")
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
	idStr := c.Param("userId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid user id")
	}

	user, err := h.usecase.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "user retrieved", user)
}

func RegisterUserHandler(group *echo.Group, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userUsecaseInstance := usecase.NewUserUsecase(userRepository)
	userHandler := NewUserHandler(userUsecaseInstance)

	userHandler.RegisterRoutes(group)
}
