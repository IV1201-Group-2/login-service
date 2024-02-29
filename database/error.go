package database

import (
	"errors"
	"fmt"
)

// Represents an error that occurred in the database layer.
type Error struct {
	// Description of the error for logging.
	Description string `json:"description"`
	// Internal wrapped error. Not visible to users.
	Internal error `json:"-"`
}

// Describes the database error.
func (e *Error) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Description, e.Internal)
	}
	return e.Description
}

// Attaches an internal error to a database error.
func (e *Error) Wrap(err error) *Error {
	return &Error{Description: e.Description, Internal: err}
}

// If an error has been wrapped in a.Internal, return the error.
func (e *Error) Unwrap() error {
	return e.Internal
}

// Database errors are considered equivalent if their description is equivalent.
func (e *Error) Is(target error) bool {
	var databaseErr *Error
	if errors.As(target, &databaseErr) {
		return e.Description == databaseErr.Description
	}
	return false
}

var (
	// ErrConnectionFailed indicates that connection to the database failed.
	ErrConnectionFailed = &Error{"connection failed", nil}
	// ErrQueryFailed indicates that an SQL query failed for an unknown reason.
	ErrQueryFailed = &Error{"query failed", nil}
	// ErrUserNotFound indicates that a user with the specificed identity couldn't be found.
	ErrUserNotFound = &Error{"user not found in db", nil}
)
