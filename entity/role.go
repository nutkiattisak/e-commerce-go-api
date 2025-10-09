package entity

import "time"

type Role struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:50;not null;uniqueIndex" json:"name"` // 'ADMIN', 'USER', 'SHOP'
	CreatedAt *time.Time `gorm:"default:now()" json:"createdAt"`
}