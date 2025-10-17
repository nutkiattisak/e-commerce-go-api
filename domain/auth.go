package domain

import (
	"context"
	"ecommerce-go-api/entity"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error)
	RegisterShop(ctx context.Context, req *entity.RegisterShopRequest) (*entity.RegisterShopResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error)
	RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error)
	VerifyToken(ctx context.Context, token string) (*entity.User, error)
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Role, error)
	AssignUserRole(ctx context.Context, userRole *entity.UserRole) error
	CreateShop(ctx context.Context, shop *entity.Shop) error
}
