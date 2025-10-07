package entity

import "time"

type Cart struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"uniqueIndex" json:"user_id"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}