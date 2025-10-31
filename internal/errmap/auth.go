package errmap

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token has expired")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)
