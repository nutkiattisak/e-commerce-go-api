package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: r}
}

func mapAddressToResponse(a *entity.Address) *entity.AddressResponse {
	if a == nil {
		return nil
	}

	subDistrictNameTh, subDistrictNameEn := "", ""
	districtNameTh, districtNameEn := "", ""
	provinceNameTh, provinceNameEn := "", ""

	if a.SubDistrict != (entity.SubDistrict{}) {
		subDistrictNameTh = a.SubDistrict.NameTH
		subDistrictNameEn = a.SubDistrict.NameEN
		if a.SubDistrict.District != nil {
			districtNameTh = a.SubDistrict.District.NameTH
			districtNameEn = a.SubDistrict.District.NameEN
			if a.SubDistrict.District.Province != nil {
				provinceNameTh = a.SubDistrict.District.Province.NameTH
				provinceNameEn = a.SubDistrict.District.Province.NameEN
			}
		}
	}

	if districtNameTh == "" && a.District != (entity.District{}) {
		districtNameTh = a.District.NameTH
		districtNameEn = a.District.NameEN
	}
	if provinceNameTh == "" && a.Province != (entity.Province{}) {
		provinceNameTh = a.Province.NameTH
		provinceNameEn = a.Province.NameEN
	}

	return &entity.AddressResponse{
		ID:                a.ID,
		UserID:            a.UserID,
		Name:              a.Name,
		Line1:             a.Line1,
		Line2:             a.Line2,
		SubDistrictID:     a.SubDistrictID,
		SubDistrictNameTh: subDistrictNameTh,
		SubDistrictNameEn: subDistrictNameEn,
		DistrictNameTh:    districtNameTh,
		DistrictNameEn:    districtNameEn,
		DistrictID:        a.DistrictID,
		ProvinceID:        a.ProvinceID,
		ProvinceNameTh:    provinceNameTh,
		ProvinceNameEn:    provinceNameEn,
		Zipcode:           a.Zipcode,
		PhoneNumber:       a.PhoneNumber,
		IsDefault:         a.IsDefault,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.UserResponse, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entity.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ImageURL:    user.ImageURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*entity.AddressResponse, error) {
	addresses, err := u.repo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.AddressResponse, 0, len(addresses))
	for _, addr := range addresses {
		result = append(result, mapAddressToResponse(addr))
	}
	return result, nil
}

func (u *userUsecase) CreateAddress(ctx context.Context, addr *entity.Address) (*entity.AddressResponse, error) {
	if addr == nil {
		return nil, nil
	}
	if err := u.repo.CreateAddress(ctx, addr); err != nil {
		return nil, err
	}

	created, err := u.repo.GetAddressByID(ctx, addr.ID)
	if err != nil {
		return mapAddressToResponse(addr), nil
	}
	return mapAddressToResponse(created), nil
}

func (u *userUsecase) GetAddressByID(ctx context.Context, id uint32) (*entity.AddressResponse, error) {
	addr, err := u.repo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapAddressToResponse(addr), nil
}

func (u *userUsecase) UpdateAddress(ctx context.Context, addr *entity.Address, userID uuid.UUID) (*entity.AddressResponse, error) {
	if err := u.repo.UpdateAddress(ctx, addr, userID); err != nil {
		return nil, err
	}

	updated, err := u.repo.GetAddressByID(ctx, addr.ID)
	if err != nil {
		return mapAddressToResponse(addr), nil
	}
	return mapAddressToResponse(updated), nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, user *entity.User) error {
	if err := u.repo.UpdateProfile(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) DeleteAddress(ctx context.Context, id uint32, userID uuid.UUID) error {
	addr, err := u.repo.GetAddressByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	if addr == nil {
		return gorm.ErrRecordNotFound
	}
	if addr.UserID != userID {
		return errmap.ErrForbidden
	}

	if addr.IsDefault {
		return errmap.ErrForbidden
	}

	return u.repo.DeleteAddress(ctx, id, userID)
}
