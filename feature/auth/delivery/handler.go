package delivery

import (
	"net/http"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/response"

	authRepo "ecommerce-go-api/feature/auth/repository"
	authUsecase "ecommerce-go-api/feature/auth/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.RegisterRequest	true	"Register payload"
//	@Success		201		{object}	entity.User
//	@Failure		400		{object}	object
//	@Router			/api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req entity.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.authUsecase.Register(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}
		return response.Error(c, statusCode, err.Error())
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", user)
}

// RegisterShop godoc
//
//	@Summary		Register a new shop
//	@Description	Create a new shop account (SHOP role)
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.RegisterShopRequest	true	"Register shop payload"
//	@Success		201		{object}	object
//	@Failure		400		{object}	object
//	@Router			/api/auth/register/shop [post]
func (h *AuthHandler) RegisterShop(c echo.Context) error {
	var req entity.RegisterShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.authUsecase.RegisterShop(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		return response.Error(c, statusCode, err.Error())
	}

	return response.Success(c, http.StatusCreated, "Shop registered successfully", resp)
}

// Login godoc
//
//	@Summary		Login
//	@Description	Authenticate user and return access & refresh tokens
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.LoginRequest	true	"Login payload"
//	@Success		200		{object}	object
//	@Failure		400		{object}	object
//	@Failure		401		{object}	object
//	@Router			/api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req entity.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	loginResp, err := h.authUsecase.Login(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid email or password" {
			statusCode = http.StatusUnauthorized
		} else if err.Error() == "user account is not active" {
			statusCode = http.StatusForbidden
		}
		return response.Error(c, statusCode, err.Error())
	}

	return response.Success(c, http.StatusOK, "Login successful", loginResp)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh access token using refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.RefreshTokenRequest	true	"Refresh payload"
//	@Success		200		{object}	object
//	@Failure		401		{object}	object
//	@Router			/api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req entity.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	loginResp, err := h.authUsecase.RefreshToken(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error())
	}

	return response.Success(c, http.StatusOK, "Token refreshed successfully", loginResp)
}

func RegisterAuthHandler(group *echo.Group, db *gorm.DB) {
	authRepository := authRepo.NewAuthRepository(db)
	authUsecaseInstance := authUsecase.NewAuthUsecase(authRepository)
	authHandler := NewAuthHandler(authUsecaseInstance)

	authHandler.RegisterRoutes(group)
}
