package entity

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NameTH    string         `gorm:"size:150;not null" json:"name_th"`
	NameEN    string         `gorm:"size:150;not null" json:"name_en"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}