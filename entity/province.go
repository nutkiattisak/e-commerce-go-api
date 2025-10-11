package entity

import (
	"time"
)

type Province struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	NameTH    string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN    string     `gorm:"size:150;not null" json:"nameEn"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
