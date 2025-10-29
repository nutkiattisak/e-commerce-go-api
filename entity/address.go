package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_addresses_user_id" json:"userId"`
	Name          string         `gorm:"type:text" json:"name"`
	Line1         string         `gorm:"type:text" json:"line1"`
	Line2         string         `gorm:"type:text" json:"line2"`
	SubDistrictID int            `gorm:"not null" json:"subDistrictId"`
	DistrictID    int            `gorm:"not null" json:"districtId"`
	ProvinceID    int            `gorm:"not null" json:"provinceId"`
	Zipcode       int            `json:"zipcode"`
	PhoneNumber   string         `gorm:"size:15" json:"phoneNumber"`
	IsDefault     bool           `gorm:"default:false;uniqueIndex:uq_addresses_user_default,where:is_default = true" json:"isDefault"`
	CreatedAt     time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"default:null" json:"deletedAt"`

	User        User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SubDistrict SubDistrict `gorm:"foreignKey:SubDistrictID;references:ID" json:"subDistrict,omitempty"`
	District    District    `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	Province    Province    `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}

type AddressResponse struct {
	ID                int       `json:"id"`
	UserID            uuid.UUID `json:"userId"`
	Name              string    `json:"name"`
	Line1             string    `json:"line1"`
	Line2             string    `json:"line2"`
	SubDistrictID     int       `json:"subDistrictId"`
	SubDistrictNameTh string    `json:"subDistrictNameTh"`
	SubDistrictNameEn string    `json:"subDistrictNameEn"`
	DistrictNameTh    string    `json:"districtNameTh"`
	DistrictNameEn    string    `json:"districtNameEn"`
	DistrictID        int       `json:"districtId"`
	ProvinceID        int       `json:"provinceId"`
	ProvinceNameTh    string    `json:"provinceNameTh"`
	ProvinceNameEn    string    `json:"provinceNameEn"`
	Zipcode           int       `json:"zipcode"`
	PhoneNumber       string    `json:"phoneNumber"`
	IsDefault         bool      `json:"isDefault"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateAddressRequest struct {
	Name          string `json:"name" validate:"required"`
	Line1         string `json:"line1" validate:"required"`
	Line2         string `json:"line2"`
	SubDistrictID int    `json:"subDistrictId" validate:"required"`
	DistrictID    int    `json:"districtId" validate:"required"`
	ProvinceID    int    `json:"provinceId" validate:"required"`
	Zipcode       int    `json:"zipcode"`
	PhoneNumber   string `json:"phoneNumber" validate:"required"`
	IsDefault     bool   `json:"isDefault"`
}

type UpdateAddressRequest struct {
	Name          string `json:"name" validate:"required,max=255"`
	Line1         string `json:"line1" validate:"required,max=500"`
	Line2         string `json:"line2" validate:"max=500"`
	SubDistrictID int    `json:"subDistrictId" validate:"required,gt=0"`
	DistrictID    int    `json:"districtId" validate:"required,gt=0"`
	ProvinceID    int    `json:"provinceId" validate:"required,gt=0"`
	Zipcode       int    `json:"zipcode" validate:"omitempty,gt=0"`
	PhoneNumber   string `json:"phoneNumber" validate:"required,max=20"`
	IsDefault     bool   `json:"isDefault"`
}
