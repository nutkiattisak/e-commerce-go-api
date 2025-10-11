package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"userId"`
	RoleID    int       `gorm:"primaryKey;not null" json:"roleId"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"-"`
}
