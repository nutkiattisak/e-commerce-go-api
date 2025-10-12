package errmap

import "errors"

var (
	ErrForbidden = errors.New("forbidden action")
	ErrNotFound  = errors.New("resource not found")
	ErrConflict  = errors.New("resource conflict")
	ErrInternal  = errors.New("internal server error")
)
