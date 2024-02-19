// The package service contains the implementation of the microservice.
// This includes handling the login process and signing JWT tokens.
package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Custom error handler conformant with shared API rules.
// https://echo.labstack.com/docs/error-handling
func ErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	var details *model.ErrorDetails
	code := http.StatusInternalServerError

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		details = &model.ErrorDetails{Message: fmt.Sprintf("%v", httpErr.Message)}
		code = httpErr.Code
	}

	if err := c.JSON(code, model.ErrorResponse{
		Error:   model.APIErrUnknown,
		Details: details,
	}); err != nil {
		c.Logger().Error(err)
	}
}

// Custom validator that uses go-playground/validator.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new instance of service.Validator.
func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

// Validates user data using go-playground/validator.
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("validation failed: %s", err.Error()))
	}

	return nil
}
