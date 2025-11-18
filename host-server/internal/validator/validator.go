package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator is a wrapper for the go-playground/validator.
type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
