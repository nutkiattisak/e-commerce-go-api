package errmap

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrProductInactive  = errors.New("product is not active")
	ErrInvalidProductID = errors.New("invalid product id")
)
