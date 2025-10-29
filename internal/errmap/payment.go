package errmap

import "errors"

var (
	ErrPaymentMethodRequired = errors.New("payment method is required")
)
