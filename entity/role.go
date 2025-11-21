package entity

const (
	RoleAdmin uint32 = 1
	RoleUser  uint32 = 2
	RoleShop  uint32 = 3
)

const (
	RoleNameAdmin = "ADMIN"
	RoleNameUser  = "USER"
	RoleNameShop  = "SHOP"
)

type Role struct {
	ID   uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null;uniqueIndex" json:"name"`
}
