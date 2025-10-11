package usecase

import (
	"context"
	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/hash"
	"ecommerce-go-api/internal/jwt"
	"errors"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type authUsecase struct {
	authRepo domain.AuthRepository
}

func NewAuthUsecase(authRepo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{
		authRepo: authRepo,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error) {

	existingUser, _ := u.authRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
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
		ImageURL:    &req.ImageURL,
	}

	if err := u.authRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	userRole, err := u.authRepo.GetRoleByName(ctx, "USER")
	if err != nil {
		return nil, err
	}

	if err := u.authRepo.AssignUserRole(ctx, &entity.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *authUsecase) RegisterShop(ctx context.Context, req *entity.RegisterShopRequest) (*entity.RegisterShopResponse, error) {

	existingUser, _ := u.authRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
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

	if err := u.authRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	shopRole, err := u.authRepo.GetRoleByName(ctx, "SHOP")
	if err != nil {
		return nil, err
	}

	if err := u.authRepo.AssignUserRole(ctx, &entity.UserRole{
		UserID: user.ID,
		RoleID: shopRole.ID,
	}); err != nil {
		return nil, err
	}

	shop := &entity.Shop{
		UserID:      user.ID,
		Name:        req.ShopName,
		Description: req.ShopDescription,
		ImageURL:    req.ShopImageURL,
		Address:     req.ShopAddress,
		IsActive:    true,
	}

	if err := u.authRepo.CreateShop(ctx, shop); err != nil {
		return nil, err
	}

	user.Password = ""

	shop.User = *user

	return &entity.RegisterShopResponse{
		Shop:         shop,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {

	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	roles, err := u.authRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	role := "USER"
	for _, r := range roles {
		if r.Name == "SHOP" {
			role = "SHOP"
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

	return &entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.LoginResponse, error) {

	claims, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := u.authRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	role := claims.Role
	if role == "" {
		role = "USER"
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

	return &entity.LoginResponse{
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
