package entity

const (
	RefundMethodBankTransfer uint32 = 1
	RefundMethodCreditCard   uint32 = 2
)

type RefundMethod struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
