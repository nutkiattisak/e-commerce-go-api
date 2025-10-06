package entity

import (
	"time"

	"gorm.io/gorm"
)

type District struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID uint           `gorm:"not null;index:idx_districts_province_id" json:"province_id"`
	NameTH     string         `gorm:"size:150;not null" json:"name_th"`
	NameEN     string         `gorm:"size:150;not null" json:"name_en"`
	CreatedAt  time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	Province Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}