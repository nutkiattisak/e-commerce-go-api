package errmap

import "errors"

var (
	ErrPaymentMethodRequired = errors.New("payment method is required")
	ErrPaymentAlreadyExists  = errors.New("payment already exists for this order")
	ErrPaymentAmountMismatch = errors.New("payment amount does not match order total")
)
