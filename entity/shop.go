package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_shops_user_id" json:"userId"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:text" json:"imageUrl"`
	Address     string         `gorm:"type:text" json:"address"`
	IsActive    bool           `gorm:"default:true;index:idx_shops_is_active" json:"isActive"`
	CreatedAt   time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"default:null" json:"-" swaggerignore:"true"`

	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

type CreateShopRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Address     string `json:"address" validate:"required"`
}

type UpdateShopRequest struct {
	Name        *string `json:"name" validate:"required,max=255"`
	Description *string `json:"description"`
	ImageURL    *string `json:"imageUrl"`
	Address     *string `json:"address"`
}

type ShopListRequest struct {
	Page       int    `query:"page"`
	PerPage    int    `query:"perPage"`
	SearchText string `query:"searchText"`
}

type ShopListResponse struct {
	Items []*Shop `json:"items"`
	Total int64   `json:"total"`
}

type ShopResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	Address     string    `json:"address"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
