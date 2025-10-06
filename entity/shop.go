package entity

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `json:"user_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	Address     string         `gorm:"type:text" json:"address"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}