package entity

import (
	"time"
)

type Courier struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	ImageURL  string     `gorm:"type:text" json:"image_url"`
	Rate      float64    `gorm:"type:decimal(10,2);not null" json:"rate"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}