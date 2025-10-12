package domain

import (
	"context"

	"ecommerce-go-api/entity"
)

// LocationRepository defines read-only accessors for administrative divisions
type LocationRepository interface {
	GetProvinces(ctx context.Context) ([]*entity.Province, error)
	GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error)
	GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error)
}

// LocationUsecase defines usecase operations for locations
type LocationUsecase interface {
	GetProvinces(ctx context.Context) ([]*entity.Province, error)
	GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error)
	GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error)
}
