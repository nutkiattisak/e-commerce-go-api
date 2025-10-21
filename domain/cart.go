package domain

import (
	"context"

	"github.com/google/uuid"

	"ecommerce-go-api/entity"
)

type CartRepository interface {
	EnsureCartForUser(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	ListCartItems(ctx context.Context, cartID int) ([]*entity.CartItem, error)
	GetCartItemsByIDs(ctx context.Context, ids []int) ([]*entity.CartItem, error)
	AddCartItem(ctx context.Context, item *entity.CartItem) error
	UpsertCartItem(ctx context.Context, item *entity.CartItem) (*entity.CartItem, bool, error)
	GetCartItemByID(ctx context.Context, id int) (*entity.CartItem, error)
	GetCartItemByCartAndProduct(ctx context.Context, cartID int, productID int) (*entity.CartItem, error)
	GetCartItemByUserAndProduct(ctx context.Context, userID uuid.UUID, productID int) (*entity.CartItem, error)
	UpdateCartItem(ctx context.Context, item *entity.CartItem) error
	DeleteCartItem(ctx context.Context, id int) error
	ClearCart(ctx context.Context, cartID int) error
}

type CartUsecase interface {
	AddItem(ctx context.Context, userID uuid.UUID, productID int, qty int) (*entity.CartItem, bool, error)
	GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, *entity.CartSummary, error)
	UpdateItem(ctx context.Context, userID uuid.UUID, itemID int, qty int) (*entity.CartItem, error)
	DeleteItem(ctx context.Context, userID uuid.UUID, itemID int) error
	EstimateShipping(ctx context.Context, userID uuid.UUID, cartItemIDs []int) (*entity.CartShippingEstimateResponse, error)
}
