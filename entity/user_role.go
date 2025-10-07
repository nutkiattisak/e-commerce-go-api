package entity

import "time"

type UserRole struct {
	UserID    int        `gorm:"primaryKey" json:"user_id"`
	RoleID    int        `gorm:"primaryKey" json:"role_id"`
	CreatedAt *time.Time `gorm:"default:now()" json:"created_at"`
	
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}
