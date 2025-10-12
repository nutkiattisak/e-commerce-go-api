package usecase

import (
	"context"

	"github.com/google/uuid"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.repo.GetByID(ctx, userID)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetByID(ctx, id)
}
