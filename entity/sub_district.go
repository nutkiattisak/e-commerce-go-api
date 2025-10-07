package entity

import (
	"time"
)

type SubDistrict struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Zipcode    int        `gorm:"not null" json:"zipcode"`
	NameTH     string     `gorm:"size:150;not null" json:"name_th"`
	NameEN     string     `gorm:"size:150;not null" json:"name_en"`
	DistrictID int        `gorm:"not null" json:"district_id"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	
	District District `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
}