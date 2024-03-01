package database_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/tests"
	// Imports ChaiSQL driver.
	_ "github.com/chaisql/chai/driver"
	"github.com/stretchr/testify/require"
)

type testError struct{}

func (t *testError) Error() string {
	return ""
}

// Tests that database.Error behaves as expected.
func TestServiceErrors(t *testing.T) {
	t.Parallel()

	wrappedError1 := errors.New("test error 1")
	wrappedError2 := errors.New("test error 2")
	wrappedError3 := errors.New("test error 3")

	error1 := database.ErrUserNotFound.Wrap(wrappedError1)
	error2 := database.ErrUserNotFound.Wrap(wrappedError2)
	error3 := database.ErrConnectionFailed.Wrap(wrappedError2)

	// errors.Is should return true on wrapped errors
	require.ErrorIs(t, error1, database.ErrUserNotFound)
	require.ErrorIs(t, error2, database.ErrUserNotFound)
	require.NotErrorIs(t, error3, database.ErrUserNotFound)
	require.NotErrorIs(t, error1, errors.New("test error 4"))

	// errors.Unwrap should return a reference to the wrapped error
	require.Equal(t, wrappedError1, errors.Unwrap(error1))
	require.NotEqual(t, wrappedError1, errors.Unwrap(error2))
	require.Equal(t, wrappedError3, errors.Unwrap(error1.Wrap(wrappedError3)))

	// errors.As should cast to database.Error correctly
	var databaseError *database.Error
	var genericError *testError

	require.ErrorAs(t, error1, &databaseError)
	require.False(t, errors.As(error1, &genericError))
}

// Test that database.Open fails on an invalid connection string.
func TestConnectionFailed(t *testing.T) {
	t.Parallel()

	_, err := database.Open(tests.RandomStr(16))
	require.ErrorIs(t, err, database.ErrConnectionFailed)
}

// Test that queries fail with an error if the connection is down.
func TestConnectionDown(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("chai", ":memory:")
	require.NoError(t, err)
	require.NoError(t, db.Close())
	require.Error(t, db.Ping())

	repository := database.NewUserRepository(db)

	_, err = repository.Query("test")
	require.ErrorIs(t, err, database.ErrQueryFailed)
}
