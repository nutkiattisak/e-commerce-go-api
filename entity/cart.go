package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"userId"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}
