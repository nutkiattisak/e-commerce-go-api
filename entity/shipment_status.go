package entity

const (
	ShipmentStatusInTransit      uint32 = 1
	ShipmentStatusDelivered      uint32 = 2
	ShipmentStatusFailedDelivery uint32 = 3
)

type ShipmentStatus struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}

type ShipmentStatusResponse struct {
	ID   uint32 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
