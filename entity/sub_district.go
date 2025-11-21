package entity

import (
	"time"
)

type SubDistrict struct {
	ID         uint32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Zipcode    uint32     `gorm:"not null" json:"zipcode"`
	NameTH     string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN     string     `gorm:"size:150;not null" json:"nameEn"`
	DistrictID uint32     `gorm:"not null" json:"districtId"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `gorm:"default:null" json:"deletedAt"`

	District *District `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
}

type SubDistrictResponse struct {
	ID         uint32 `json:"id"`
	Zipcode    uint32 `json:"zipcode"`
	NameTH     string `json:"nameTh"`
	NameEN     string `json:"nameEn"`
	DistrictID uint32 `json:"districtId"`
}
