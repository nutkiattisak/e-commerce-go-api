package jwt

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token has expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := getJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
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
	duration := os.Getenv("JWT_ACCESS_TOKEN_DURATION")
	if duration == "" {
		return 15 * time.Minute
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

func getRefreshTokenDuration() time.Duration {
	duration := os.Getenv("JWT_REFRESH_TOKEN_DURATION")
	if duration == "" {
		return 7 * 24 * time.Hour // 7 days
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 7 * 24 * time.Hour
	}
	return d
}
