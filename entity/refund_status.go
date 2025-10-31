package entity

const (
	RefundStatusPending   int = 1
	RefundStatusApproved  int = 2
	RefundStatusCompleted int = 3
	RefundStatusRejected  int = 4
)

type RefundStatus struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
