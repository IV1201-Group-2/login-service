package service

import (
	"errors"
	"fmt"
)

// Represents an error that occurred in the service layer.
type Error struct {
	// Description of the error for logging.
	Description string `json:"description"`
	// Internal wrapped error. Not visible to users.
	Internal error `json:"-"`
}

// Describes the service error.
func (e *Error) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Description, e.Internal)
	}
	return e.Description
}

// Attaches an internal error to a service error.
func (e *Error) Wrap(err error) *Error {
	return &Error{Description: e.Description, Internal: err}
}

// If an error has been wrapped in a.Internal, return the error.
func (e *Error) Unwrap() error {
	return e.Internal
}

// Service errors are considered equivalent if their description is equivalent.
func (e *Error) Is(target error) bool {
	var databaseErr *Error
	if errors.As(target, &databaseErr) {
		return e.Description == databaseErr.Description
	}
	return false
}

var (
	// ErrWrongIdentity indicates that authentication failed because account with that identity was not found.
	ErrWrongIdentity = &Error{"wrong identity", nil}
	// ErrWrongPassword indicates that authentication failed because the wrong password was provided.
	ErrWrongPassword = &Error{"wrong password", nil}
	// ErrMissingPassword indicates that authentication failed because the user has no password in the database.
	ErrMissingPassword = &Error{"missing password", nil}
	// ErrBcryptError indicates that password update failed because Bcrypt returned an error.
	ErrBcryptError = &Error{"bcrypt error", nil}
	// ErrJWTError indicates that authentication failed because golang-jwt returned an error.
	ErrJWTError = &Error{"jwt error", nil}
)
