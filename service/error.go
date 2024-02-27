package service

import (
	"fmt"
	"net/http"
)

// Represents an error that occurred in the service layer.
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

// Attaches detailed user-visible information to an API error.
// This is intended to give the API consumer more information about where and how it occurred.
func (a *Error) WithDetails(details any) *Error {
	return &Error{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: details, Internal: a.Internal}
}

// Attaches an internal error to an API error.
func (a *Error) WithInternal(err error) *Error {
	return &Error{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: a.Details, Internal: err}
}

// Describes the API error.
func (a *Error) Error() string {
	if a.Internal != nil {
		return fmt.Sprintf("%s: %v", a.ErrorType, a.Internal)
	}
	return a.ErrorType
}

// If an error has been wrapped in a.Internal, return the error.
func (a *Error) Unwrap() error {
	return a.Internal
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

	// ErrWrongIdentity indicates that no account was found with that specific username or email address.
	ErrWrongIdentity = &Error{http.StatusUnauthorized, "WRONG_IDENTITY", nil, nil}
	// ErrWrongPassword indicates that an account was found but the wrong password was provided.
	ErrWrongPassword = &Error{http.StatusUnauthorized, "WRONG_PASSWORD", nil, nil}

	// ErrAlreadyLoggedIn indicates that the user is already logged in (JWT token was provided).
	ErrAlreadyLoggedIn = &Error{http.StatusBadRequest, "ALREADY_LOGGED_IN", nil, nil}
	// ErrTokenNotProvided indicates that the user did not provide a token for reset API.
	ErrTokenNotProvided = &Error{http.StatusUnauthorized, "TOKEN_NOT_PROVIDED", nil, nil} // #nosec G101
	// ErrTokenInvalid indicates that the user provided an invalid or expired token.
	ErrTokenInvalid = &Error{http.StatusUnauthorized, "INVALID_TOKEN", nil, nil}

	// ErrInvalidRoute indicates that the user tried to access an invalid route.
	ErrInvalidRoute = &Error{http.StatusNotFound, "INVALID_ROUTE", nil, nil}
)
