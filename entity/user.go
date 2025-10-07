package entity

import "time"

type User struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string     `gorm:"size:255;not null" json:"first_name"`
	LastName    string     `gorm:"size:255;not null" json:"last_name"`
	Email       string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password    string     `gorm:"type:text;not null" json:"-"`
	PhoneNumber string     `gorm:"size:15;not null;index:idx_users_phone_number" json:"phone_number"`
	ImageURL    string     `gorm:"type:text" json:"image_url"`
	CreatedAt   *time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"default:now()" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	Shops []Shop `gorm:"foreignKey:UserID" json:"shops,omitempty"`
}