package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint           `json:"user_id"`
	Name          string         `gorm:"type:text" json:"name"`
	Line1         string         `gorm:"type:text" json:"line1"`
	Line2         string         `gorm:"type:text" json:"line2"`
	SubDistrictID uint           `gorm:"not null" json:"sub_district_id"`
	DistrictID    uint           `gorm:"not null" json:"district_id"`
	ProvinceID    uint           `gorm:"not null" json:"province_id"`
	Zipcode       int            `json:"zipcode"`
	PhoneNumber   string         `gorm:"size:15" json:"phone_number"`
	IsDefault     bool           `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	User        User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SubDistrict SubDistrict `gorm:"foreignKey:SubDistrictID;references:ID" json:"sub_district,omitempty"`
	District    District    `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	Province    Province    `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}