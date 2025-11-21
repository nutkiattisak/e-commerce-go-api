package entity

const (
	PaymentStatusPending    uint32 = 1
	PaymentStatusProcessing uint32 = 2
	PaymentStatusCompleted  uint32 = 3
	PaymentStatusFailed     uint32 = 4
	PaymentStatusCancelled  uint32 = 5
	PaymentStatusRefunded   uint32 = 6
	PaymentStatusExpired    uint32 = 7
)

type PaymentStatus struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
