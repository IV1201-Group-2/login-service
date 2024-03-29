package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// Tests that the server returns a valid JWT token when a user logs in.
func TestLogin(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": tests.MockPassword,
		"role":     tests.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	obj := model.LoginTokenResponse{}
	body, _ := io.ReadAll(res.Body)

	// Parse the response
	require.NoError(t, json.Unmarshal(body, &obj))
	require.NotEqual(t, "", obj.Token, "Response does not contain token")

	claims := model.UserClaims{}
	// Parse the embedded JWT token
	_, err := jwt.ParseWithClaims(obj.Token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.Equal(t, tests.MockApplicant.Email, claims.Email)
	require.Equal(t, tests.MockApplicant.Role, claims.Role)
	require.Equal(t, "login", claims.Usage)
}

// Tests that the server returns MISSING_PARAMETERS when API caller is missing required parameters.
func TestLoginMissingParameters(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"password": tests.MockPassword,
		"role":     tests.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "MISSING_PARAMETERS", obj.ErrorType)
}

// Tests that the server does not return MISSING_PARAMETERS when API caller is missing optional parameters.
func TestLoginOptionalParameters(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": tests.MockPassword,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}

// Tests that the server returns WRONG_IDENTITY when user does not exist.
func TestLoginMissingUser(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": "doesnotexist@example.com",
		"password": "password",
		"role":     1,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_IDENTITY", obj.ErrorType)
}

// Tests that the server returns WRONG_IDENTITY when user has a different role.
func TestLoginWrongRole(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": tests.MockPassword,
		"role":     model.RoleRecruiter,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_IDENTITY", obj.ErrorType)
}

// Tests that the server returns WRONG_IDENTITY when user has wrong password.
func TestLoginWrongPassword(t *testing.T) {
	t.Parallel()

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": "wrong",
		"role":     model.RoleApplicant,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_IDENTITY", obj.ErrorType)
}

// Tests that the server returns MISSING_PASSWORD and a reset token when user has wrong password.
func TestLoginMissingPassword(t *testing.T) {
	t.Parallel()

	// Try the mock applicant that won't be modified during tests
	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant2.Email,
		"password": tests.MockPassword,
		"role":     model.RoleApplicant,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusNotFound, res.StatusCode)

	details := model.ResetTokenResponse{}
	obj := api.Error{Details: &details}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "MISSING_PASSWORD", obj.ErrorType)

	require.NotEqual(t, "", details.Token, "Response does not contain token")

	claims := model.UserClaims{}
	// Parse the embedded JWT token
	_, err := jwt.ParseWithClaims(details.Token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.Equal(t, tests.MockApplicant2.Email, claims.Email)
	require.Equal(t, "reset", claims.Usage)
}

// Tests that the server returns ALREADY_LOGGED_IN when a JWT token is set.
func TestAlreadyLoggedIn(t *testing.T) {
	t.Parallel()

	testToken, _, _ := service.SignUserToken(tests.MockApplicant, []byte(os.Getenv("JWT_SECRET")))

	res := tests.Request(t, "/api/login", map[string]any{
		"identity": tests.MockApplicant.Email,
		"password": tests.MockPassword,
		"role":     tests.MockApplicant.Role,
	}, map[string]string{
		"Authorization": "Bearer " + testToken,
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	obj := api.Error{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "ALREADY_LOGGED_IN", obj.ErrorType)
}
