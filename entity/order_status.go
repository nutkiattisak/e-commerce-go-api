package entity

const (
	OrderStatusPending    int = 1
	OrderStatusProcessing int = 2
	OrderStatusShipped    int = 3
	OrderStatusDelivered  int = 4
	OrderStatusCompleted  int = 5
	OrderStatusCancelled  int = 6
)

type OrderStatus struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
