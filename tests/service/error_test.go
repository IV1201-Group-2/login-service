package service_test

import (
	"errors"
	"testing"

	"github.com/IV1201-Group-2/login-service/service"
	"github.com/stretchr/testify/require"
)

type testError struct{}

func (t *testError) Error() string {
	return ""
}

// Tests that service.Error behaves as expected.
func TestServiceErrors(t *testing.T) {
	t.Parallel()

	wrappedError1 := errors.New("test error 1")
	wrappedError2 := errors.New("test error 2")
	wrappedError3 := errors.New("test error 3")

	error1 := service.ErrJWTError.Wrap(wrappedError1)
	error2 := service.ErrJWTError.Wrap(wrappedError2)
	error3 := service.ErrBcryptError.Wrap(wrappedError2)

	// errors.Is should return true on wrapped errors
	require.ErrorIs(t, error1, service.ErrJWTError)
	require.ErrorIs(t, error2, service.ErrJWTError)
	require.NotErrorIs(t, error3, service.ErrJWTError)
	require.NotErrorIs(t, error1, errors.New("test error 4"))

	// errors.Unwrap should return a reference to the wrapped error
	require.Equal(t, wrappedError1, errors.Unwrap(error1))
	require.NotEqual(t, wrappedError1, errors.Unwrap(error2))
	require.Equal(t, wrappedError3, errors.Unwrap(error1.Wrap(wrappedError3)))

	// errors.As should cast to service.Error correctly
	var serviceError *service.Error
	var genericError *testError

	require.ErrorAs(t, error1, &serviceError)
	require.False(t, errors.As(error1, &genericError))
}
