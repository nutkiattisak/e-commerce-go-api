package entity

import (
	"time"

	"gorm.io/gorm"
)

type SubDistrict struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Zipcode    int            `gorm:"not null" json:"zipcode"`
	NameTH     string         `gorm:"size:150;not null" json:"name_th"`
	NameEN     string         `gorm:"size:150;not null" json:"name_en"`
	DistrictID uint           `gorm:"not null;index:idx_sub_districts_district_id" json:"district_id"`
	CreatedAt  time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	// Relations
	District District `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
}