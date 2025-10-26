package delivery

import (
	"errors"
	"net/http"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
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
//	@Success		201		{object}	entity.RegisterResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req entity.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.authUsecase.Register(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errmap.ErrEmailAlreadyExists) {
			statusCode = http.StatusConflict
		}
		return response.Error(c, statusCode, errmap.ErrEmailAlreadyExists.Error())
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
//	@Success		201		{object}	entity.RegisterShopResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/auth/register-shop [post]
func (h *AuthHandler) RegisterShop(c echo.Context) error {
	var req entity.RegisterShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.authUsecase.RegisterShop(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errmap.ErrEmailAlreadyExists) {
			statusCode = http.StatusConflict
		}
		return response.Error(c, statusCode, errmap.ErrEmailAlreadyExists.Error())
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
//	@Success		200		{object}	entity.AuthResponse
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req entity.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	loginResponse, err := h.authUsecase.Login(c.Request().Context(), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errmap.ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
		}
		return response.Error(c, statusCode, errmap.ErrInvalidCredentials.Error())
	}

	return response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh access token using refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		entity.RefreshTokenRequest	true	"Refresh payload"
//	@Success		200		{object}	entity.AuthResponse
//	@Failure		400		{object}	response.ResponseError
//	@Failure		401		{object}	response.ResponseError
//	@Failure		500		{object}	response.ResponseError
//	@Router			/api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req entity.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, errmap.ErrInvalidRequest.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	loginResponse, err := h.authUsecase.RefreshToken(c.Request().Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, errmap.ErrInvalidRefreshToken):
			return response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			c.Logger().Error("RefreshToken error: ", err)
			return response.Error(c, http.StatusInternalServerError, errmap.ErrInternalServer.Error())
		}
	}

	return response.Success(c, http.StatusOK, "Token refreshed successfully", loginResponse)
}

func RegisterAuthHandler(group *echo.Group, db *gorm.DB) {
	authRepository := authRepo.NewAuthRepository(db)
	authUsecaseInstance := authUsecase.NewAuthUsecase(authRepository)
	authHandler := NewAuthHandler(authUsecaseInstance)
	authHandler.RegisterRoutes(group)
}
