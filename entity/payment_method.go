package entity

const (
	PaymentMethodCreditCard   uint32 = 1
	PaymentMethodCod          uint32 = 2
	PaymentMethodBankTransfer uint32 = 3
	PaymentMethodPromptPay    uint32 = 4
)

type PaymentMethod struct {
	ID   uint32 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
