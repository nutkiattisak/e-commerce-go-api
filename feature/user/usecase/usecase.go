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

func (u *userUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*entity.Address, error) {
	return u.repo.GetAddressesByUserID(ctx, userID)
}

func (u *userUsecase) CreateAddress(ctx context.Context, addr *entity.Address) (*entity.Address, error) {
	if addr == nil {
		return nil, nil
	}
	if err := u.repo.CreateAddress(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (u *userUsecase) GetAddressByID(ctx context.Context, id int) (*entity.Address, error) {
	return u.repo.GetAddressByID(ctx, id)
}

func (u *userUsecase) UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) (*entity.Address, error) {
	if err := u.repo.UpdateAddress(ctx, addr, userID); err != nil {
		return nil, err
	}
	return addr, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := u.repo.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) DeleteAddress(ctx context.Context, id int, userID uuid.UUID) error {
	return u.repo.DeleteAddress(ctx, id, userID)
}
