package entity

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex" json:"name"` // 'ADMIN', 'USER', 'SELLER'
	CreatedAt time.Time `json:"created_at"`
}