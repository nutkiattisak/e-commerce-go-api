package entity

import (
	"time"

	"gorm.io/gorm"
)

// UserType constants
const (
	UserTypeBuyer  = "buyer"
	UserTypeSeller = "seller"
)

type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string         `gorm:"size:255" json:"first_name"`
	LastName    string         `gorm:"size:255" json:"last_name"`
	Email       string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"type:text;not null" json:"-"`
	PhoneNumber string         `gorm:"size:15;not null;index:idx_users_phone_number" json:"phone_number"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	UserType    string         `gorm:"size:20;not null;default:'buyer'" json:"user_type"` // "buyer" or "seller"
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}