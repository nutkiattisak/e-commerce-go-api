package validator

import (
	v "github.com/go-playground/validator/v10"
)

type Validator struct {
	V *v.Validate
}

func New() *Validator {
	return &Validator{V: v.New()}
}

// Echo expects a Validate(i interface{}) error method
func (vld *Validator) Validate(i interface{}) error {
	return vld.V.Struct(i)
}
