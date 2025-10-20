package errmap

import "errors"

var (
	ErrForbidden      = errors.New("forbidden")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrConflict       = errors.New("conflict")
)
