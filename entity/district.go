package entity

import (
	"time"
)

type District struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID int        `gorm:"not null" json:"provinceId"`
	NameTH     string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN     string     `gorm:"size:150;not null" json:"nameEn"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`

	Province Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
}
