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

func mapToProvinceResponse(p *entity.Province) *entity.ProvinceResponse {
	if p == nil {
		return nil
	}
	return &entity.ProvinceResponse{
		ID:     p.ID,
		NameTH: p.NameTH,
		NameEN: p.NameEN,
	}
}

func mapToDistrictResponse(d *entity.District) *entity.DistrictResponse {
	if d == nil {
		return nil
	}
	resp := &entity.DistrictResponse{
		ID:         d.ID,
		ProvinceID: d.ProvinceID,
		NameTH:     d.NameTH,
		NameEN:     d.NameEN,
	}
	if d.Province != nil {
		resp.Province = mapToProvinceResponse(d.Province)
	}
	return resp
}

func mapToSubDistrictResponse(s *entity.SubDistrict) *entity.SubDistrictResponse {
	if s == nil {
		return nil
	}
	return &entity.SubDistrictResponse{
		ID:         s.ID,
		Zipcode:    s.Zipcode,
		NameTH:     s.NameTH,
		NameEN:     s.NameEN,
		DistrictID: s.DistrictID,
	}
}

func (u *locationUsecase) GetProvinces(ctx context.Context) ([]*entity.ProvinceResponse, error) {
	provinces, err := u.repo.GetProvinces(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*entity.ProvinceResponse, len(provinces))
	for i, p := range provinces {
		responses[i] = mapToProvinceResponse(p)
	}

	return responses, nil
}

func (u *locationUsecase) GetDistrictsByProvince(ctx context.Context, provinceID int) ([]*entity.DistrictResponse, error) {
	districts, err := u.repo.GetDistrictsByProvince(ctx, provinceID)
	if err != nil {
		return nil, err
	}

	responses := make([]*entity.DistrictResponse, len(districts))
	for i, d := range districts {
		responses[i] = mapToDistrictResponse(d)
	}

	return responses, nil
}

func (u *locationUsecase) GetSubDistrictsByDistrict(ctx context.Context, districtID int) ([]*entity.SubDistrictResponse, error) {
	subs, err := u.repo.GetSubDistrictsByDistrict(ctx, districtID)
	if err != nil {
		return nil, err
	}

	responses := make([]*entity.SubDistrictResponse, len(subs))
	for i, s := range subs {
		responses[i] = mapToSubDistrictResponse(s)
	}

	return responses, nil
}
