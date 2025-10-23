package domain

import (
	"context"

	"ecommerce-go-api/entity"
)

type CourierUsecase interface {
	ListCouriers(ctx context.Context) ([]entity.CourierListResponse, error)
}

type CourierRepository interface {
	ListAll(ctx context.Context) ([]*entity.Courier, error)
}
