package model

import (
	"fmt"
	"net/http"
)

// Represents an error that occurred in the service layer.
type APIError struct {
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
func (a *APIError) WithDetails(details any) *APIError {
	return &APIError{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: details, Internal: a.Internal}
}

// Attaches an internal error to an API error.
func (a *APIError) WithInternal(err error) *APIError {
	return &APIError{StatusCode: a.StatusCode, ErrorType: a.ErrorType, Details: a.Details, Internal: err}
}

// Describes the API error.
func (a *APIError) Error() string {
	if a.Internal != nil {
		return fmt.Sprintf("%s (%v)", a.ErrorType, a.Internal)
	}
	return a.ErrorType
}

// If an error has been wrapped in a.Internal, return the error.
func (a *APIError) Unwrap() error {
	return a.Internal
}

var (
	// Unknown error.
	ErrUnknown = &APIError{http.StatusInternalServerError, "UNKNOWN", nil, nil}
	// An external service such as the database is unavailable.
	ErrServiceUnavailable = &APIError{http.StatusInternalServerError, "SERVICE_UNAVAILABLE", nil, nil}

	// User did not provide identity, password or desired role.
	ErrMissingParameters = &APIError{http.StatusBadRequest, "MISSING_PARAMETERS", nil, nil}
	// User does not have a password in the database.
	ErrMissingPassword = &APIError{http.StatusNotFound, "MISSING_PASSWORD", nil, nil}

	// No account was found with that specific username or email address.
	ErrWrongIdentity = &APIError{http.StatusUnauthorized, "WRONG_IDENTITY", nil, nil}
	// Account was found but the wrong password was provided.
	ErrWrongPassword = &APIError{http.StatusUnauthorized, "WRONG_PASSWORD", nil, nil}

	// User is already logged in (JWT token was provided).
	ErrAlreadyLoggedIn = &APIError{http.StatusBadRequest, "ALREADY_LOGGED_IN", nil, nil}
	// User did not provide a token for reset API.
	ErrTokenNotProvided = &APIError{http.StatusUnauthorized, "TOKEN_NOT_PROVIDED", nil, nil} // #nosec G101
	// User provided an invalid or expired token.
	ErrTokenInvalid = &APIError{http.StatusUnauthorized, "INVALID_TOKEN", nil, nil}

	// User tried to access an invalid route.
	ErrInvalidRoute = &APIError{http.StatusNotFound, "INVALID_ROUTE", nil, nil}
)
