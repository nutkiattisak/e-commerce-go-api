package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderNumber         string     `gorm:"size:20;not null;uniqueIndex" json:"orderNumber"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;index:idx_orders_user_id" json:"userId"`
	AddressID           int        `json:"addressId"`
	GrandTotal          float64    `gorm:"type:decimal(10,2);not null" json:"grandTotal"`
	CancelReason        string     `gorm:"type:text" json:"cancelReason"`
	ShippingName        string     `gorm:"size:255;not null" json:"shippingName"`
	ShippingPhone       string     `gorm:"size:15;not null" json:"shippingPhone"`
	ShippingLine1       string     `gorm:"type:text;not null" json:"shippingLine1"`
	ShippingLine2       string     `gorm:"type:text" json:"shippingLine2"`
	ShippingSubDistrict string     `gorm:"size:100;not null" json:"shippingSubDistrict"`
	ShippingDistrict    string     `gorm:"size:100;not null" json:"shippingDistrict"`
	ShippingProvince    string     `gorm:"size:100;not null" json:"shippingProvince"`
	ShippingZipcode     string     `gorm:"size:5;not null" json:"shippingZipcode"`
	CreatedAt           time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CancelledAt         *time.Time `json:"cancelledAt"`

	User    User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Address *Address `gorm:"foreignKey:AddressID;references:ID" json:"address,omitempty"`
}
