package middleware

import (
	"ecommerce-go-api/internal/jwt"
	"ecommerce-go-api/internal/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, "Missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			}

			tokenString := parts[1]

			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					return response.Error(c, http.StatusUnauthorized, "Token has expired")
				}
				return response.Error(c, http.StatusUnauthorized, "Invalid token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RoleAuth(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			if userRole == nil {
				return response.Error(c, http.StatusUnauthorized, "Unauthorized")
			}

			role := userRole.(string)
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}

			return response.Error(c, http.StatusForbidden, "Access denied: insufficient permissions")
		}
	}
}

func AdminOnly() echo.MiddlewareFunc {
	return RoleAuth("ADMIN")
}

func ShopOwnerOnly() echo.MiddlewareFunc {
	return RoleAuth("SHOP")
}

func UserOnly() echo.MiddlewareFunc {
	return RoleAuth("USER")
}

func ShopOrAdmin() echo.MiddlewareFunc {
	return RoleAuth("SHOP", "ADMIN")
}
