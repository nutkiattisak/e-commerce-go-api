package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName   string     `gorm:"type:varchar(255);not null" json:"firstName"`
	LastName    string     `gorm:"type:varchar(255);not null" json:"lastName"`
	Email       string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password    string     `gorm:"type:text;not null" json:"-"`
	PhoneNumber string     `gorm:"type:varchar(15);not null;index" json:"phoneNumber"`
	ImageURL    *string    `gorm:"type:text" json:"imageUrl,omitempty"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	ImageURL    string `json:"imageURL"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
