package entity

const (
	PaymentStatusPending    int = 1
	PaymentStatusProcessing int = 2
	PaymentStatusCompleted  int = 3
	PaymentStatusFailed     int = 4
	PaymentStatusCancelled  int = 5
	PaymentStatusRefunded   int = 6
	PaymentStatusExpired    int = 7
)

type PaymentStatus struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
