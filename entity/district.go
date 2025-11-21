package entity

import (
	"time"

	"gorm.io/gorm"
)

type District struct {
	ID         uint32         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID uint32         `gorm:"not null" json:"provinceId"`
	NameTH     string         `gorm:"size:150;not null" json:"nameTh"`
	NameEN     string         `gorm:"size:150;not null" json:"nameEn"`
	CreatedAt  time.Time      `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"default:null" json:"deletedAt"`

	Province *Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}

type DistrictResponse struct {
	ID         uint32            `json:"id"`
	ProvinceID uint32            `json:"provinceId"`
	NameTH     string            `json:"nameTh"`
	NameEN     string            `json:"nameEn"`
	Province   *ProvinceResponse `json:"province,omitempty"`
}
