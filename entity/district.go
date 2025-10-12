package entity

import (
	"time"
)

type District struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"provinceId"`
	NameTH     string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN     string     `gorm:"size:150;not null" json:"nameEn"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`

	Province *Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}
