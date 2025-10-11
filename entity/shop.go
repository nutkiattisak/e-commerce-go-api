package entity

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_shops_user_id" json:"userId"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ImageURL    string     `gorm:"type:text" json:"imageUrl"`
	Address     string     `gorm:"type:text" json:"address"`
	IsActive    bool       `gorm:"default:true;index:idx_shops_is_active" json:"isActive"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

type CreateShopRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Address     string `json:"address" validate:"required"`
}

type UpdateShopRequest struct {
	Name        string `json:"name" validate:"omitempty,max=255"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Address     string `json:"address"`
}

type ShopListRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Search   string `query:"search"`
	SortBy   string `query:"sort_by"`
	Order    string `query:"order"`
}

type ShopListResponse struct {
	Shops      []*Shop `json:"shops"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	TotalPages int     `json:"totalPages"`
}
