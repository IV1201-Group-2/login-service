package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"

	"github.com/labstack/echo/v4"
)

// Represents an error that occurred in the API layer.
// API errors can be translated into a JSON object following shared API rules.
type Error struct {
	// HTTP status code.
	StatusCode int `json:"-"`
	// Error type.
	ErrorType string `json:"error"`
	// Error details.
	Details any `json:"details,omitempty"`
	// Internal wrapped error. Not visible to users.
	Internal error `json:"-"`
}

// Describes the API error.
func (a *Error) Error() string {
	if a.Internal != nil {
		return fmt.Sprintf("%s: %v", a.ErrorType, a.Internal)
	}
	return a.ErrorType
}

// Attaches detailed user-visible information to an API error.
// This is intended to give the API consumer more information about where and how it occurred.
func (a *Error) WithDetails(details any) *Error {
	return &Error{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: details, Internal: a.Internal}
}

// Attaches an internal error to an API error.
func (a *Error) Wrap(err error) *Error {
	return &Error{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: a.Details, Internal: err}
}

// If an error has been wrapped in a.Internal, return the error.
func (a *Error) Unwrap() error {
	return a.Internal
}

// Service errors are considered equivalent if their error type is equivalent.
func (a *Error) Is(target error) bool {
	var apiError *Error
	if errors.As(target, &apiError) {
		return a.ErrorType == apiError.ErrorType
	}
	return false
}

var (
	// ErrUnknown represents an unknown error.
	ErrUnknown = &Error{http.StatusInternalServerError, "UNKNOWN", nil, nil}
	// ErrServiceUnavailable indicates that an external service such as the database is unavailable.
	ErrServiceUnavailable = &Error{http.StatusInternalServerError, "SERVICE_UNAVAILABLE", nil, nil}

	// ErrMissingParameters indicates that the user did not provide identity, password or desired role.
	ErrMissingParameters = &Error{http.StatusBadRequest, "MISSING_PARAMETERS", nil, nil}
	// ErrMissingParameters indicates that the user does not have a password in the database.
	ErrMissingPassword = &Error{http.StatusNotFound, "MISSING_PASSWORD", nil, nil}

	// ErrWrongIdentity indicates that no account was found with the provided parameters.
	ErrWrongIdentity = &Error{http.StatusUnauthorized, "WRONG_IDENTITY", nil, nil}

	// ErrAlreadyLoggedIn indicates that the user is already logged in (JWT token was provided).
	ErrAlreadyLoggedIn = &Error{http.StatusBadRequest, "ALREADY_LOGGED_IN", nil, nil}
	// ErrTokenNotProvided indicates that the user did not provide a token for reset API.
	ErrTokenNotProvided = &Error{http.StatusUnauthorized, "TOKEN_NOT_PROVIDED", nil, nil} // #nosec G101
	// ErrTokenInvalid indicates that the user provided an invalid or expired token.
	ErrTokenInvalid = &Error{http.StatusUnauthorized, "INVALID_TOKEN", nil, nil}

	// ErrInvalidRoute indicates that the user tried to access an invalid route.
	ErrInvalidRoute = &Error{http.StatusNotFound, "INVALID_ROUTE", nil, nil}
)

// Rewrites errors returned by Echo to follow shared API rules.
func rewriteEchoErrors(err *echo.HTTPError, c echo.Context) *Error {
	result := ErrUnknown.Wrap(err)
	if errors.Is(echo.ErrNotFound, err) ||
		errors.Is(echo.ErrForbidden, err) ||
		errors.Is(echo.ErrMethodNotAllowed, err) {
		result = ErrInvalidRoute.Wrap(err)
	}
	logging.Errorf(c, "Rewrote framework error %d to %s", err.Code, result.ErrorType)
	return result
}

// Custom error handler conformant with shared API rules.
func ErrorHandler(err error, c echo.Context) {
	userVisibleErr := ErrUnknown.Wrap(err)

	var apiErr *Error
	var databaseErr *database.Error
	var httpErr *echo.HTTPError

	switch {
	case errors.As(err, &apiErr):
		if internalErr := apiErr.Unwrap(); internalErr != nil {
			logging.Errorf(c, "Error occurred in handler: %v", internalErr)
		}
		userVisibleErr = apiErr
	case errors.As(err, &httpErr):
		logging.Errorf(c, "Error occurred in framework: %v", err)
		// Special case for some framework errors
		userVisibleErr = rewriteEchoErrors(httpErr, c)
	case errors.As(err, &databaseErr):
		logging.Errorf(c, "Error occurred in database: %v", err)
		userVisibleErr = ErrServiceUnavailable.Wrap(err)
	default:
		logging.Errorf(c, "Recovered from unexpected error: %v", err)
	}

	if err := c.JSON(userVisibleErr.StatusCode, userVisibleErr); err != nil {
		logging.Errorf(c, "Error occurred in HTTP error handler: %v", err)
	}
}
