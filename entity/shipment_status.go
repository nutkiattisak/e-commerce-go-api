package entity

const (
	ShipmentStatusPending        int = 1
	ShipmentStatusReadyToShip    int = 2
	ShipmentStatusInTransit      int = 3
	ShipmentStatusDelivered      int = 4
	ShipmentStatusFailedDelivery int = 5
	ShipmentStatusReturned       int = 6
	ShipmentStatusCancelled      int = 7
)

type ShipmentStatus struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
