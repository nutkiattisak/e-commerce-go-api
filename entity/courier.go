package entity

import (
	"time"

	"gorm.io/gorm"
)

type Courier struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	ImageURL  string         `gorm:"type:text" json:"image_url"`
	Rate      float64        `gorm:"type:decimal(10,2);not null" json:"rate"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}