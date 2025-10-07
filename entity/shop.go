package entity

import (
	"time"
)

type Shop struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int        `gorm:"not null;index:idx_shops_user_id" json:"user_id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ImageURL    string     `gorm:"type:text" json:"image_url"`
	Address     string     `gorm:"type:text" json:"address"`
	IsActive    bool       `gorm:"default:true;index:idx_shops_is_active" json:"is_active"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}