package errmap

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
