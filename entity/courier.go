package entity

import (
	"time"
)

type Courier struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	ImageURL  string     `gorm:"type:text" json:"imageUrl"`
	Rate      float64    `gorm:"type:decimal(10,2);not null" json:"rate"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
