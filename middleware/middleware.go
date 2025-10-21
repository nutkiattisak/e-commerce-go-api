package middleware

import (
	"net/http"
	"strings"

	"ecommerce-go-api/internal/jwt"
	"ecommerce-go-api/internal/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	CTX_KEY_USER_ID = "userId"
	CTX_KEY_ROLE    = "role"
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
				if err == jwt.ErrExpiredToken {
					return response.Error(c, http.StatusUnauthorized, "Token has expired")
				}
				if err == jwt.ErrInvalidSigningMethod {
					return response.Error(c, http.StatusUnauthorized, "Invalid token signature")
				}
				return response.Error(c, http.StatusUnauthorized, "Invalid token")
			}

			if claims.UserID == uuid.Nil {
				return response.Error(c, http.StatusUnauthorized, "Invalid user ID")
			}

			c.Set(CTX_KEY_USER_ID, claims.UserID)
			c.Set(CTX_KEY_ROLE, claims.Role)

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
			roleKey := c.Get(CTX_KEY_ROLE)
			if roleKey == nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized")
			}

			role, ok := roleKey.(string)
			if !ok || role == "" {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized")
			}

			if _, exists := roleSet[role]; !exists {
				return response.Error(c, http.StatusForbidden, "Access denied: insufficient permissions")
			}

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, error) {
	v := c.Get("userId")
	if v == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	switch t := v.(type) {
	case uuid.UUID:
		return t, nil
	case string:
		id, err := uuid.Parse(t)
		if err != nil {
			return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}
		return id, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid user id type")
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
