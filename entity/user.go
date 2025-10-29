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
	DeletedAt   gorm.DeletedAt `gorm:"default:null" json:"deletedAt"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email" example:"kiattisak.c@example.com"`
	Password    string `json:"password" validate:"required,min=8" example:"yourpassword"`
	FirstName   string `json:"firstName" validate:"required" example:"Kiattisak"`
	LastName    string `json:"lastName" validate:"required" example:"Chantharamaneechote"`
	PhoneNumber string `json:"phoneNumber" validate:"required" example:"0900000000"`
	ImageURL    string `json:"imageUrl" example:"https://example.com/image.jpg"`
}

type RegisterResponse struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	ImageURL    *string   `json:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RegisterShopRequest struct {
	Email           string `json:"email" validate:"required,email" example:"kiattisak.c@example.com"`
	Password        string `json:"password" validate:"required,min=8" example:"yourpassword"`
	FirstName       string `json:"firstName" validate:"required" example:"Kiattisak"`
	LastName        string `json:"lastName" validate:"required" example:"Chantharamaneechote"`
	PhoneNumber     string `json:"phoneNumber" validate:"required" example:"0900000000"`
	ImageURL        string `json:"imageUrl"`
	ShopName        string `json:"shopName" validate:"required,max=255" example:"Nike Shop"`
	ShopDescription string `json:"shopDescription" example:"The best Nike products available."`
	ShopImageURL    string `json:"shopImageUrl"`
	ShopAddress     string `json:"shopAddress" validate:"required" example:"321 Moo 4 Suthep Road, Chiang Mai"`
}

type RegisterShopResponse struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ImageURL    string    `gorm:"type:text" json:"imageUrl"`
	Address     string    `gorm:"type:text" json:"address"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	User *RegisterResponse `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"kiattisak.c@example.com"`
	Password string `json:"password" validate:"required" example:"yourpassword"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" example:"your_refresh_token"`
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
	FirstName   string  `json:"firstName" validate:"required" example:"Kiattisak"`
	LastName    string  `json:"lastName" validate:"required" example:"Chantharamaneechote"`
	PhoneNumber string  `json:"phoneNumber" validate:"required" example:"0900000000"`
	ImageURL    *string `json:"imageUrl" example:"https://example.com/image.jpg"`
}
