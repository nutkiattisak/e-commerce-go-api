package entity

import (
	"time"
)

type Address struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int        `gorm:"not null;index:idx_addresses_user_id" json:"user_id"`
	Name          string     `gorm:"type:text" json:"name"`
	Line1         string     `gorm:"type:text" json:"line1"`
	Line2         string     `gorm:"type:text" json:"line2"`
	SubDistrictID int        `gorm:"not null" json:"sub_district_id"`
	DistrictID    int        `gorm:"not null" json:"district_id"`
	ProvinceID    int        `gorm:"not null" json:"province_id"`
	Zipcode       int        `json:"zipcode"`
	PhoneNumber   string     `gorm:"size:15" json:"phone_number"`
	IsDefault     bool       `gorm:"default:false;uniqueIndex:uq_addresses_user_default,where:is_default = true" json:"is_default"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	
	User        User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SubDistrict SubDistrict `gorm:"foreignKey:SubDistrictID;references:ID" json:"sub_district,omitempty"`
	District    District    `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	Province    Province    `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}