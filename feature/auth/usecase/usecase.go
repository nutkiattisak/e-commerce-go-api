package usecase

import (
	"context"
	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
	"ecommerce-go-api/internal/hash"
	"ecommerce-go-api/internal/jwt"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type authUsecase struct {
	authRepo domain.AuthRepository
}

func NewAuthUsecase(authRepo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{
		authRepo: authRepo,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	existingUser, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errmap.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	if req.ImageURL != "" {
		user.ImageURL = &req.ImageURL
	}

	if err := u.authRepo.RegisterUser(ctx, user); err != nil {
		return nil, err
	}

	return &entity.RegisterResponse{
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

func (u *authUsecase) RegisterShop(ctx context.Context, req *entity.RegisterShopRequest) (*entity.RegisterShopResponse, error) {
	existingUser, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errmap.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	if req.ImageURL != "" {
		user.ImageURL = &req.ImageURL
	}

	shop := &entity.Shop{
		Name:        req.ShopName,
		Description: req.ShopDescription,
		ImageURL:    req.ShopImageURL,
		Address:     req.ShopAddress,
		IsActive:    true,
	}

	if err := u.authRepo.RegisterShop(ctx, user, shop); err != nil {
		return nil, err
	}

	user.Password = ""
	shop.User = user

	return &entity.RegisterShopResponse{
		ID:          shop.ID,
		Name:        shop.Name,
		Description: shop.Description,
		ImageURL:    shop.ImageURL,
		Address:     shop.Address,
		CreatedAt:   shop.CreatedAt,
		UpdatedAt:   shop.UpdatedAt,
		User: &entity.RegisterResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			ImageURL:    user.ImageURL,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error) {

	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errmap.ErrNotFound) {
			return nil, errmap.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, errmap.ErrInvalidCredentials
	}

	roles, err := u.authRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	role := entity.RoleNameUser
	for _, r := range roles {
		if r.ID == entity.RoleShop {
			role = entity.RoleNameShop
			break
		}
		if r.ID == entity.RoleAdmin {
			role = entity.RoleNameAdmin
			break
		}
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &entity.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error) {

	claims, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errmap.ErrInvalidRefreshToken
	}

	user, err := u.authRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	roles, err := u.authRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for refresh: %w", err)
	}

	role := entity.RoleNameUser
	for _, r := range roles {
		if r.ID == entity.RoleShop {
			role = entity.RoleNameShop
			break
		}
		if r.ID == entity.RoleAdmin {
			role = entity.RoleNameAdmin
			break
		}
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &entity.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *authUsecase) VerifyToken(ctx context.Context, token string) (*entity.User, error) {

	claims, err := jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := u.authRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
