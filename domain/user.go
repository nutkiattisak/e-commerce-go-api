package domain

import (
	"context"

	"github.com/google/uuid"

	"ecommerce-go-api/entity"
)

type UserUsecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]*entity.AddressResponse, error)
	CreateAddress(ctx context.Context, addr *entity.Address) (*entity.AddressResponse, error)
	GetAddressByID(ctx context.Context, id int) (*entity.AddressResponse, error)
	UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) (*entity.AddressResponse, error)
	UpdateProfile(ctx context.Context, user *entity.User) error
	DeleteAddress(ctx context.Context, id int, userID uuid.UUID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetAddressesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Address, error)
	CreateAddress(ctx context.Context, addr *entity.Address) error
	GetAddressByID(ctx context.Context, id int) (*entity.Address, error)
	UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) error
	UpdateProfile(ctx context.Context, user *entity.User) error
	DeleteAddress(ctx context.Context, id int, userID uuid.UUID) error
}
