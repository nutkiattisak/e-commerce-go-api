package errmap

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidUserIDType = errors.New("invalid user id type")
)
