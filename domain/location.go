package domain

import (
	"context"

	"ecommerce-go-api/entity"
)

type LocationRepository interface {
	GetProvinces(ctx context.Context) ([]*entity.Province, error)
	GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error)
	GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error)
}

type LocationUsecase interface {
	GetProvinces(ctx context.Context) ([]*entity.Province, error)
	GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error)
	GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error)
}
