package jwt

import (
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/timeth"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	secret := getJWTSecret()
	duration := getAccessTokenDuration()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeth.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(timeth.Now()),
			NotBefore: jwt.NewNumericDate(timeth.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	secret := getJWTSecret()
	duration := getRefreshTokenDuration()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeth.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(timeth.Now()),
			NotBefore: jwt.NewNumericDate(timeth.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := getJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errmap.ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errmap.ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(timeth.Now()) {
		return nil, errmap.ErrExpiredToken
	}

	return claims, nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is not set")
	}
	return secret
}

func getAccessTokenDuration() time.Duration {
	const defaultDuration = 15 * time.Minute // 15 minutes

	durationStr := os.Getenv("JWT_ACCESS_TOKEN_DURATION")

	if durationStr == "" {
		return defaultDuration
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf(
			"Warning: Invalid JWT_ACCESS_TOKEN_DURATION format '%s'. Using default 15m.",
			durationStr,
		)
		return defaultDuration
	}

	return d
}

func getRefreshTokenDuration() time.Duration {
	const defaultRefreshDuration = 7 * 24 * time.Hour // 7 days

	durationStr := os.Getenv("JWT_REFRESH_TOKEN_DURATION")

	if durationStr == "" {
		return defaultRefreshDuration
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf(
			"Warning: Invalid JWT_REFRESH_TOKEN_DURATION format '%s'. Using default 7 days.",
			durationStr,
		)
		return defaultRefreshDuration
	}

	return d
}
