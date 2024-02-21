// The package service contains the implementation of the microservice.
// This includes handling the login process and signing JWT tokens.
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Log an error that occured in a handler.
func LogErrorf(c echo.Context, format string, args ...any) {
	c.Logger().Errorf(fmt.Sprintf("[%s] %s", c.RealIP(), format), args...)
}

// Rewrites errors returned by Echo to follow shared API rules.
func rewriteEchoErrors(err *echo.HTTPError, c echo.Context) *model.APIError {
	result := model.ErrUnknown.WithInternal(err)
	if errors.Is(echo.ErrNotFound, err) ||
		errors.Is(echo.ErrForbidden, err) ||
		errors.Is(echo.ErrMethodNotAllowed, err) {
		result = model.ErrInvalidRoute.WithInternal(err)
	}
	LogErrorf(c, "Rewrote framework error %d to %s", err.Code, result.ErrorType)
	return result
}

// Custom error handler conformant with shared API rules.
func ErrorHandler(err error, c echo.Context) {
	userVisibleErr := model.ErrUnknown.WithInternal(err)

	var apiErr *model.APIError
	var httpErr *echo.HTTPError

	switch {
	case errors.As(err, &apiErr):
		if internalErr := apiErr.Unwrap(); internalErr != nil {
			LogErrorf(c, "Error occurred in handler: %v", internalErr)
		}
		userVisibleErr = apiErr
	case errors.As(err, &httpErr):
		LogErrorf(c, "Error occurred in framework: %v", err)
		// Special case for some framework errors
		userVisibleErr = rewriteEchoErrors(httpErr, c)
	case errors.Is(err, sql.ErrConnDone):
	case errors.Is(err, sql.ErrTxDone):
	case errors.Is(err, net.ErrClosed):
	case errors.Is(err, syscall.ECONNREFUSED):
	case errors.Is(err, syscall.ECONNABORTED):
	case errors.Is(err, syscall.ECONNRESET):
		LogErrorf(c, "Error occurred in database: %v", err)
		userVisibleErr = model.ErrServiceUnavailable.WithInternal(err)
	default:
		LogErrorf(c, "Recovered from unexpected error: %v", err)
	}

	if err := c.JSON(userVisibleErr.StatusCode, userVisibleErr); err != nil {
		LogErrorf(c, "Error occurred in HTTP error handler: %v", err)
	}
}

// Custom validator that uses go-playground/validator.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new instance of service.Validator.
func NewValidator() *Validator {
	return &Validator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// Validates user data using go-playground/validator.
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return model.ErrMissingParameters.WithInternal(err)
	}

	return nil
}
