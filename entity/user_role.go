package entity

import "time"

type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}