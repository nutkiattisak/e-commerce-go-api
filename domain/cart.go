package domain

import (
	"context"

	"github.com/google/uuid"

	"ecommerce-go-api/entity"
)

type CartUsecase interface {
	AddItem(ctx context.Context, userID uuid.UUID, productID uint32, qty uint32) (*entity.CartItem, bool, error)
	GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, *entity.CartSummary, error)
	UpdateItem(ctx context.Context, userID uuid.UUID, itemID uint32, qty uint32) (*entity.CartItem, error)
	DeleteItem(ctx context.Context, userID uuid.UUID, itemID uint32) error
	EstimateShipping(ctx context.Context, userID uuid.UUID, cartItemIDs []uint32) (*entity.CartShippingEstimateResponse, error)
}

type CartRepository interface {
	EnsureCartForUser(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	ListCartItems(ctx context.Context, cartID uint32) ([]*entity.CartItem, error)
	GetCartItemsByIDs(ctx context.Context, ids []uint32) ([]*entity.CartItem, error)
	AddCartItem(ctx context.Context, item *entity.CartItem) error
	UpsertCartItem(ctx context.Context, item *entity.CartItem) (*entity.CartItem, bool, error)
	GetCartItemByID(ctx context.Context, id uint32) (*entity.CartItem, error)
	GetCartItemByCartAndProduct(ctx context.Context, cartID uint32, productID uint32) (*entity.CartItem, error)
	GetCartItemByUserAndProduct(ctx context.Context, userID uuid.UUID, productID uint32) (*entity.CartItem, error)
	UpdateCartItem(ctx context.Context, item *entity.CartItem) error
	DeleteCartItem(ctx context.Context, id uint32) error
	ClearCart(ctx context.Context, cartID uint32) error
}
