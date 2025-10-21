package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName   string         `gorm:"type:varchar(255);not null" json:"firstName"`
	LastName    string         `gorm:"type:varchar(255);not null" json:"lastName"`
	Email       string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password    string         `gorm:"type:text;not null" json:"-"`
	PhoneNumber string         `gorm:"type:varchar(15);not null;index" json:"phoneNumber"`
	ImageURL    *string        `gorm:"type:text" json:"imageUrl,omitempty"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"default:now()" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"default:null" json:"deletedAt" swaggerignore:"true"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	ImageURL    string `json:"imageUrl"`
}

type RegisterShopRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	ImageURL        string `json:"imageUrl"`
	ShopName        string `json:"shopName" validate:"required,max=255"`
	ShopDescription string `json:"shopDescription"`
	ShopImageURL    string `json:"shopImageUrl"`
	ShopAddress     string `json:"shopAddress" validate:"required"`
}

type RegisterShopResponse struct {
	Shop *Shop `json:"shop"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateProfileRequest struct {
	FirstName   string  `json:"firstName" validate:"required"`
	LastName    string  `json:"lastName" validate:"required"`
	PhoneNumber string  `json:"phoneNumber" validate:"required"`
	ImageURL    *string `json:"imageUrl"`
}
