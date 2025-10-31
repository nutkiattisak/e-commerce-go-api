package entity

const (
	PaymentMethodCreditCard   int = 1
	PaymentMethodCod          int = 2
	PaymentMethodBankTransfer int = 3
	PaymentMethodPromptPay    int = 4
)

type PaymentMethod struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Name string `gorm:"size:100;not null" json:"name"`
}
