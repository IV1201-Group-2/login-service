package api_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/stretchr/testify/require"
)

// Tests that api.Error.Is and api.Error.Unwrap behaves as expected
func TestAPIErrors(t *testing.T) {
	t.Parallel()

	wrappedError1 := errors.New("test error 1")
	wrappedError2 := errors.New("test error 2")
	wrappedError3 := errors.New("test error 3")

	error1 := api.ErrWrongIdentity.Wrap(wrappedError1)
	error2 := api.ErrWrongIdentity.Wrap(wrappedError2)
	error3 := api.ErrMissingParameters.Wrap(wrappedError2)

	// errors.Is should return true on wrapped errors
	require.ErrorIs(t, error1, api.ErrWrongIdentity)
	require.ErrorIs(t, error2, api.ErrWrongIdentity)
	require.NotErrorIs(t, error3, api.ErrWrongIdentity)

	// errors.Unwrap should return a reference to the wrapped error
	require.Equal(t, errors.Unwrap(error1), wrappedError1)
	require.NotEqual(t, errors.Unwrap(error2), wrappedError1)
	require.Equal(t, errors.Unwrap(error1.Wrap(wrappedError3)), wrappedError3)

	// errors.As should cast to api.Error correctly
	var apiError *api.Error
	require.ErrorAs(t, error1, &apiError)
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
