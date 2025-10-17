package entity

import (
	"time"
)

type SubDistrict struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Zipcode    int        `gorm:"not null" json:"zipcode"`
	NameTH     string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN     string     `gorm:"size:150;not null" json:"nameEn"`
	DistrictID int        `gorm:"not null" json:"districtId"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `gorm:"default:null" json:"deletedAt"`

	District *District `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
}

type SubDistrictResponse struct {
	ID         int    `json:"id"`
	Zipcode    int    `json:"zipcode"`
	NameTH     string `json:"nameTh"`
	NameEN     string `json:"nameEn"`
	DistrictID int    `json:"districtId"`
}
