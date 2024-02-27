package api

import (
	"github.com/go-playground/validator/v10"
)

// Custom validator that uses go-playground/validator.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new instance of service.Validator.
func NewValidator() *Validator {
	return &Validator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// Validates user data using go-playground/validator.
func (cv *Validator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return ErrMissingParameters.Wrap(err)
	}

	return nil
}
