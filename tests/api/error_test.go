package api_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/tests"
	// Imports ChaiSQL driver.
	_ "github.com/chaisql/chai/driver"
	"github.com/stretchr/testify/require"
)

type testError struct{}

func (t *testError) Error() string {
	return ""
}

// Tests that api.Error behaves as expected.
func TestAPIErrors(t *testing.T) {
	t.Parallel()

	wrappedError1 := errors.New("test error 1")
	wrappedError2 := errors.New("test error 2")
	wrappedError3 := errors.New("test error 3")

	error1 := api.ErrWrongIdentity.Wrap(wrappedError1)
	error2 := api.ErrWrongIdentity.Wrap(wrappedError2)
	error3 := api.ErrMissingParameters.Wrap(wrappedError2)
	require.NotErrorIs(t, error1, errors.New("test error 4"))

	// errors.Is should return true on wrapped errors
	require.ErrorIs(t, error1, api.ErrWrongIdentity)
	require.ErrorIs(t, error2, api.ErrWrongIdentity)
	require.NotErrorIs(t, error3, api.ErrWrongIdentity)

	// errors.Unwrap should return a reference to the wrapped error
	require.Equal(t, wrappedError1, errors.Unwrap(error1))
	require.NotEqual(t, wrappedError1, errors.Unwrap(error2))
	require.Equal(t, wrappedError3, errors.Unwrap(error1.Wrap(wrappedError3)))

	// errors.As should cast to api.Error correctly
	var apiError *api.Error
	var genericError *testError

	require.ErrorAs(t, error1, &apiError)
	require.False(t, errors.As(error1, &genericError))
}

// Tests that the server returns an error conformant with shared API rules on wrong route.
func TestWrongRoute(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/wrong", map[string]any{}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusNotFound, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "INVALID_ROUTE", obj.ErrorType)
}

// Tests that the server can handle an invalid JWT token.
func TestInvalidToken(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/wrong", map[string]any{}, map[string]string{
		"Authorization": "Bearer " + tests.RandomStr(16),
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "INVALID_TOKEN", obj.ErrorType)
}

// Tests that the server can handle the database connection going down gracefully.
func TestDatabaseConnectionDown(t *testing.T) {
	t.Parallel()

	// Open a temporary database connection.
	db, err := sql.Open("chai", ":memory:")
	require.NoError(t, err)

	srv, err := api.NewServer(db)
	require.NoError(t, err)
	defer srv.Close()

	// Close the connection and make sure it's down.
	require.NoError(t, db.Close())
	require.Error(t, db.Ping())

	res := tests.CustomRequest(t, srv, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": tests.MockPassword,
		"role":     tests.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "SERVICE_UNAVAILABLE", obj.ErrorType)
}
