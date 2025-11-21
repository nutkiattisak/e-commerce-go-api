package entity

import (
	"time"

	"gorm.io/gorm"
)

type Courier struct {
	ID        uint32         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	ImageURL  string         `gorm:"type:text" json:"imageUrl"`
	Rate      float64        `gorm:"type:decimal(10,2);not null" json:"rate"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"default:null" json:"deletedAt"`
}
type CourierListResponse struct {
	ID       uint32  `json:"id"`
	Name     string  `json:"name"`
	ImageURL string  `json:"imageUrl"`
	Rate     float64 `json:"rate"`
}
