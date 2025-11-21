package entity

const (
	OrderStatusPending    uint32 = 1
	OrderStatusProcessing uint32 = 2
	OrderStatusShipped    uint32 = 3
	OrderStatusDelivered  uint32 = 4
	OrderStatusCompleted  uint32 = 5
	OrderStatusCancelled  uint32 = 6
)

type OrderStatus struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
