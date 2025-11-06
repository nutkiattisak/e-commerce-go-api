package entity

const (
	RoleAdmin int = 1
	RoleUser  int = 2
	RoleShop  int = 3
)

const (
	RoleNameAdmin = "ADMIN"
	RoleNameUser  = "USER"
	RoleNameShop  = "SHOP"
)

type Role struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null;uniqueIndex" json:"name"`
}
