package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"userId"`
	RoleID    int       `gorm:"primaryKey;not null" json:"roleId"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"-"`
}
