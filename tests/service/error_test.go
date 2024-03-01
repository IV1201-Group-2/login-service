package service_test

import (
	"errors"
	"testing"

	"github.com/IV1201-Group-2/login-service/service"
	"github.com/stretchr/testify/require"
)

// Tests that service.Error behaves as expected
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

	// errors.Unwrap should return a reference to the wrapped error
	require.Equal(t, errors.Unwrap(error1), wrappedError1)
	require.NotEqual(t, errors.Unwrap(error2), wrappedError1)
	require.Equal(t, errors.Unwrap(error1.Wrap(wrappedError3)), wrappedError3)

	// errors.As should cast to service.Error correctly
	var apiError *service.Error
	require.ErrorAs(t, error1, &apiError)
}
