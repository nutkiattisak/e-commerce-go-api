package entity

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index:idx_addresses_user_id" json:"userId"`
	Name          string     `gorm:"type:text" json:"name"`
	Line1         string     `gorm:"type:text" json:"line1"`
	Line2         string     `gorm:"type:text" json:"line2"`
	SubDistrictID int        `gorm:"not null" json:"subDistrictId"`
	DistrictID    int        `gorm:"not null" json:"districtId"`
	ProvinceID    int        `gorm:"not null" json:"provinceId"`
	Zipcode       int        `json:"zipcode"`
	PhoneNumber   string     `gorm:"size:15" json:"phoneNumber"`
	IsDefault     bool       `gorm:"default:false;uniqueIndex:uq_addresses_user_default,where:is_default = true" json:"isDefault"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
	
	User        User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SubDistrict SubDistrict `gorm:"foreignKey:SubDistrictID;references:ID" json:"subDistrict,omitempty"`
	District    District    `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	Province    Province    `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}