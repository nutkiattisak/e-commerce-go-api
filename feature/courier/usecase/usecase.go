package usecase

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type courierUsecase struct {
	repo domain.CourierRepository
}

func NewCourierUsecase(repo domain.CourierRepository) domain.CourierUsecase {
	return &courierUsecase{repo: repo}
}

func (u *courierUsecase) ListCouriers(ctx context.Context) ([]entity.CourierListResponse, error) {
	couriers, err := u.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []entity.CourierListResponse
	for _, c := range couriers {
		response = append(response, entity.CourierListResponse{
			ID:       c.ID,
			Name:     c.Name,
			ImageURL: c.ImageURL,
			Rate:     c.Rate,
		})
	}

	return response, nil
}
