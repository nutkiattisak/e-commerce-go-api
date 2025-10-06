package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRefunded   OrderStatus = "REFUNDED"
)

type Order struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderNumber           string     `gorm:"size:20;not null;uniqueIndex" json:"order_number"`
	UserID                uint       `gorm:"not null;index:idx_orders_user_id" json:"user_id"`
	AddressID             *uint      `json:"address_id"`
	GrandTotal            float64    `gorm:"type:decimal(10,2);not null" json:"grand_total"`
	CancelReason          string     `gorm:"type:text" json:"cancel_reason"`
	ShippingName          string     `gorm:"size:255;not null" json:"shipping_name"`
	ShippingPhone         string     `gorm:"size:15;not null" json:"shipping_phone"`
	ShippingLine1         string     `gorm:"type:text;not null" json:"shipping_line1"`
	ShippingLine2         string     `gorm:"type:text" json:"shipping_line2"`
	ShippingSubDistrict   string     `gorm:"size:100;not null" json:"shipping_sub_district"`
	ShippingDistrict      string     `gorm:"size:100;not null" json:"shipping_district"`
	ShippingProvince      string     `gorm:"size:100;not null" json:"shipping_province"`
	ShippingZipcode       string     `gorm:"size:5;not null" json:"shipping_zipcode"`
	CreatedAt             time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	CancelledAt           *time.Time `json:"cancelled_at"`
	
	User    User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Address *Address `gorm:"foreignKey:AddressID;references:ID" json:"address,omitempty"`
}