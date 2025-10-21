package errmap

import "errors"

var (
	ErrQuantityMustBeGreaterThanZero = errors.New("quantity must be greater than 0")
	ErrFailedToGetCartItem           = errors.New("failed to get cart item")
	ErrInsufficientStock             = errors.New("insufficient stock")
	ErrNoShippingOptions             = errors.New("no shipping options for shop")
	ErrFailedToGetUserCart           = errors.New("failed to get user's cart")
	ErrInvalidQuantity               = errors.New("qty must be greater than zero")
	ErrCartIsEmpty                   = errors.New("cart is empty")
	ErrCartItemNotFound              = errors.New("cart item not found")
)
