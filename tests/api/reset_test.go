package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// Tests that the server can reset the password of a user.
func TestPasswordReset(t *testing.T) {
	t.Parallel()

	repo := database.NewUserRepository(tests.Database)

	// Generate a new random password every time the test is run.
	newPassword := tests.RandomStr(16)
	// Create a new reset token for the user
	resetToken, _, _ := service.SignResetToken(tests.MockApplicant3, []byte(os.Getenv("JWT_SECRET")))

	// Go down into service layer and make sure we can't authenticate as this user before reset
	_, err := service.AuthenticateUser(repo, tests.MockApplicant3.Email, newPassword, nil)
	require.ErrorIs(t, err, service.ErrMissingPassword)

	// Send the request
	res := tests.Request(t, "/api/reset", map[string]any{
		"password": newPassword,
	}, map[string]string{
		"Authorization": "Bearer " + resetToken,
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	obj := model.LoginTokenResponse{}
	body, _ := io.ReadAll(res.Body)

	// Parse the response
	require.NoError(t, json.Unmarshal(body, &obj))
	require.NotEqual(t, "", obj.Token, "Response does not contain token")

	claims := model.UserClaims{}
	// Parse the embedded JWT token
	_, err = jwt.ParseWithClaims(obj.Token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.Equal(t, tests.MockApplicant3.Email, claims.Email)
	require.Equal(t, tests.MockApplicant3.Role, claims.Role)
	require.Equal(t, "login", claims.Usage)

	// Go down into service layer again and make sure we can authenticate as this user after reset
	_, err = service.AuthenticateUser(repo, tests.MockApplicant3.Email, newPassword, nil)
	require.NoError(t, err)
}

// Tests that the server returns MISSING_PARAMETERS when API caller is missing required parameters.
func TestResetMissingParameters(t *testing.T) {
	t.Parallel()

	resetToken, _, _ := service.SignResetToken(tests.MockApplicant3, []byte(os.Getenv("JWT_SECRET")))
	res := tests.Request(t, "/api/reset", map[string]any{}, map[string]string{
		"Authorization": "Bearer " + resetToken,
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "MISSING_PARAMETERS", obj.ErrorType)
}

// Test that reset functionality rejects a request without a token.
func TestResetMissingToken(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/reset", map[string]any{
		"password": tests.MockPassword,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "TOKEN_NOT_PROVIDED", obj.ErrorType)
}

// Test that reset functionality rejects a token intended for login.
func TestResetLoginToken(t *testing.T) {
	t.Parallel()

	testToken, _, _ := service.SignUserToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))

	res := tests.Request(t, "/api/reset", map[string]any{
		"password": tests.MockPassword,
	}, map[string]string{
		"Authorization": "Bearer " + testToken,
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "INVALID_TOKEN", obj.ErrorType)
}
