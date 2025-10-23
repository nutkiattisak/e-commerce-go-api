package errmap

import "errors"

var (
	ErrAddressIDRequired = errors.New("addressId is required")
)
