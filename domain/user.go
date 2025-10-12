package domain

import (
	"context"

	"github.com/google/uuid"

	"ecommerce-go-api/entity"
)

type UserUsecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
