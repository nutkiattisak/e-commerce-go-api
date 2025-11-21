package entity

const (
	RefundStatusPending   uint32 = 1
	RefundStatusApproved  uint32 = 2
	RefundStatusCompleted uint32 = 3
	RefundStatusRejected  uint32 = 4
)

type RefundStatus struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
