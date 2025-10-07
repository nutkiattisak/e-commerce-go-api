package entity

import (
	"time"
)

type District struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"province_id"`
	NameTH     string     `gorm:"size:150;not null" json:"name_th"`
	NameEN     string     `gorm:"size:150;not null" json:"name_en"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	
	Province Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}