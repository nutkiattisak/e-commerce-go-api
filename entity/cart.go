package entity

import "time"

type Cart struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}