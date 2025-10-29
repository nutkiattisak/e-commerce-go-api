package entity

type Role struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null;uniqueIndex" json:"name"`
}
