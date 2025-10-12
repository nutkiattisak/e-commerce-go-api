package usecase

import (
	"context"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
)

type locationUsecase struct {
	repo domain.LocationRepository
}

func NewLocationUsecase(r domain.LocationRepository) domain.LocationUsecase {
	return &locationUsecase{repo: r}
}

func (u *locationUsecase) GetProvinces(ctx context.Context) ([]*entity.Province, error) {
	provinces, err := u.repo.GetProvinces(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (u *locationUsecase) GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.District, error) {
	districts, err := u.repo.GetDistrictsByProvince(ctx, provinceID)
	if err != nil {
		return nil, err
	}
	return districts, nil
}

func (u *locationUsecase) GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrict, error) {
	subs, err := u.repo.GetSubDistrictsByDistrict(ctx, districtID)
	if err != nil {
		return nil, err
	}
	return subs, nil
}
