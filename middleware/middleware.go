package middleware

import (
	"errors"
	"net/http"
	"strings"

	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/jwt"
	"ecommerce-go-api/internal/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, "Missing authorization header")
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			}

			tokenString := parts[1]

			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {
				if errors.Is(err, errmap.ErrExpiredToken) {
					return response.Error(c, http.StatusUnauthorized, "Token has expired")
				}
				if errors.Is(err, errmap.ErrInvalidSigningMethod) {
					return response.Error(c, http.StatusUnauthorized, "Invalid token signature")
				}
				return response.Error(c, http.StatusUnauthorized, "Invalid token")
			}

			if claims.UserID == uuid.Nil {
				return response.Error(c, http.StatusUnauthorized, "Invalid user ID")
			}

			c.Set("userId", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RoleAuth(allowedRoles ...string) echo.MiddlewareFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleKey := c.Get("role")
			if roleKey == nil {
				return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
			}

			role, ok := roleKey.(string)
			if !ok || role == "" {
				return response.Error(c, http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
			}

			if _, exists := roleSet[role]; !exists {
				return response.Error(c, http.StatusForbidden, errmap.ErrForbidden.Error())
			}

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, error) {
	v := c.Get("userId")
	if v == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, errmap.ErrUnauthorized.Error())
	}
	switch t := v.(type) {
	case uuid.UUID:
		return t, nil
	case string:
		id, err := uuid.Parse(t)
		if err != nil {
			return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, errmap.ErrInvalidUserID.Error())
		}
		return id, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, errmap.ErrInvalidUserIDType.Error())
	}
}

func ShopOwnerOnly() echo.MiddlewareFunc {
	return RoleAuth("SHOP")
}

func UserOnly() echo.MiddlewareFunc {
	return RoleAuth("USER")
}

func ShopOrUser() echo.MiddlewareFunc {
	return RoleAuth("SHOP", "USER")
}
