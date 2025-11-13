package entity

const (
	ShipmentStatusInTransit      int = 1
	ShipmentStatusDelivered      int = 2
	ShipmentStatusFailedDelivery int = 3
)

type ShipmentStatus struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}

type ShipmentStatusResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
