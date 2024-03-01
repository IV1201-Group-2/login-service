package service_test

import (
	"os"
	"testing"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// Tests that signing a login token works correctly and all parameters are present.
func TestSignLoginToken(t *testing.T) {
	t.Parallel()

	token, expiry, err := service.SignUserToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)

	claims := model.UserClaims{}
	decodedToken, err := jwt.ParseWithClaims(token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.True(t, decodedToken.Valid)
	require.Equal(t, claims.ExpiresAt.Time, expiry)

	require.Equal(t, tests.MockApplicant.Email, claims.Email)
	require.Equal(t, tests.MockApplicant.Role, claims.Role)
	require.Equal(t, "login", claims.Usage)
}

// Tests that signing a reset token works correctly and all parameters are present.
func TestSignResetToken(t *testing.T) {
	t.Parallel()

	token, expiry, err := service.SignResetToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)

	claims := model.UserClaims{}
	decodedToken, err := jwt.ParseWithClaims(token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.True(t, decodedToken.Valid)
	require.Equal(t, claims.ExpiresAt.Time, expiry)

	require.Equal(t, tests.MockApplicant.Email, claims.Email)
	require.Equal(t, tests.MockApplicant.Role, claims.Role)
	require.Equal(t, "reset", claims.Usage)
}

// Tests that all tokens are signed with HS256.
func TestSignHS256(t *testing.T) {
	t.Parallel()

	token1, _, err := service.SignUserToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)
	token2, _, err := service.SignUserToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))
	require.NoError(t, err)

	claims := model.UserClaims{}

	decodedToken1, err := jwt.ParseWithClaims(token1, &claims, mockKeyFunc)
	require.NoError(t, err, jwt.ErrSignatureInvalid)
	require.Equal(t, decodedToken1.Method.Alg(), "HS256")

	decodedToken2, err := jwt.ParseWithClaims(token2, &claims, mockKeyFunc)
	require.NoError(t, err, jwt.ErrSignatureInvalid)
	require.Equal(t, decodedToken2.Method.Alg(), "HS256")
}

// Tests that signing a token with the wrong secret fails to verify.
func TestSignWrongSecret(t *testing.T) {
	t.Parallel()

	token, _, err := service.SignUserToken(tests.MockApplicant, []byte("wrong"))
	require.NoError(t, err)

	claims := model.UserClaims{}
	decodedToken, err := jwt.ParseWithClaims(token, &claims, mockKeyFunc)

	require.ErrorIs(t, err, jwt.ErrSignatureInvalid)
	require.False(t, decodedToken.Valid)
}
