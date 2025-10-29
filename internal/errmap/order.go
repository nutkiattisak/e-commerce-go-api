package errmap

import "errors"

var (
	ErrInvalidOrderID      = errors.New("invalid order id")
	ErrFailedToCreateOrder = errors.New("failed to create order")
	ErrFailedToGetOrder    = errors.New("failed to get order")
	ErrFailedToListOrders  = errors.New("failed to list orders")
	ErrOrderNotFound       = errors.New("order not found")
)
