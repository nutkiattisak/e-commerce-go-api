package entity

import (
	"time"
)

type Province struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	NameTH    string     `gorm:"size:150;not null" json:"nameTh"`
	NameEN    string     `gorm:"size:150;not null" json:"nameEn"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt,omitempty"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
