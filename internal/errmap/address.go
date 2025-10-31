package errmap

import "errors"

var (
	ErrAddressIDRequired = errors.New("addressId is required")
	ErrAddressNotFound   = errors.New("address not found")
	ErrInvalidAddressID  = errors.New("invalid address id")
)
