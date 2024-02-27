package database

import (
	"fmt"
)

// Represents an error that occurred in the database layer.
type Error struct {
	// Description of the error for logging.
	Description string `json:"description"`
	// Internal wrapped error. Not visible to users.
	Internal error `json:"-"`
}

// Attaches an internal error to a database error.
func (d *Error) WithInternal(err error) *Error {
	return &Error{Description: d.Description, Internal: err}
}

// Describes the database error.
func (a *Error) Error() string {
	if a.Internal != nil {
		return fmt.Sprintf("%s: %v", a.Description, a.Internal)
	}
	return a.Description
}

// If an error has been wrapped in a.Internal, return the error.
func (a *Error) Unwrap() error {
	return a.Internal
}

var (
	// ErrConnectionFailed indicates that connection to the database failed.
	ErrConnectionFailed = &Error{"connection failed", nil}

	// ErrConnectionMockMode indicates that a mock database is being used.
	// This is a warning and can be ignored if the user is informed.
	ErrConnectionMockMode = &Error{"database is in mock mode", nil}

	// ErrUserNotFound indicates that a user with the specificed identity couldn't be found.
	ErrUserNotFound = &Error{"user not found in db", nil}
)
