package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator is a wrapper for the go-playground/validator.
type Validator struct {
	validator *validator.Validate
}

// defaultValidator is the global singleton instance
var defaultValidator = &Validator{
	validator: validator.New(),
}

// New returns the global validator instance (singleton pattern)
func New() *Validator {
	return defaultValidator
}

// Validate validates a struct using the global validator instance
// This is a convenient function for direct validation without creating an instance
func Validate(i any) error {
	return defaultValidator.Validate(i)
}

// Validate validates a struct
func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return formatError(err)
	}
	return nil
}

// FormatError formats validator.ValidationErrors into a readable error message
// Returns the first validation error with a formatted message
func formatError(err error) error {
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return fmt.Errorf("%s is required", e.Field())
			case "gt":
				return fmt.Errorf("%s must be greater than %s", e.Field(), e.Param())
			case "gte":
				return fmt.Errorf("%s must be greater than or equal to %s", e.Field(), e.Param())
			case "lt":
				return fmt.Errorf("%s must be less than %s", e.Field(), e.Param())
			case "lte":
				return fmt.Errorf("%s must be less than or equal to %s", e.Field(), e.Param())
			case "oneof":
				return fmt.Errorf("%s must be one of: %s", e.Field(), e.Param())
			case "email":
				return fmt.Errorf("%s must be a valid email address", e.Field())
			case "url":
				return fmt.Errorf("%s must be a valid URL", e.Field())
			default:
				return fmt.Errorf("%s failed validation: %s", e.Field(), e.Tag())
			}
		}
	}

	return err
}
