package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

// newMockServer creates a new Echo server that uses the mock dataset.
func newMockServer() *echo.Echo {
	srv := echo.New()
	db, _ := database.Connect("mock")

	srv.HTTPErrorHandler = service.ErrorHandler
	srv.Validator = service.NewValidator()

	// TODO: Maybe better to pass a new auth config instead...
	os.Setenv("JWT_SECRET", "mocksecret")
	service.RegisterRoutes(srv, db)

	return srv
}

func mockKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte("mocksecret"), nil
}

var useJSON = os.Getenv("TEST_JSON") == "1"

// Sends a request to mock server and returns response.
func testRequest(t *testing.T, path string, params map[string]any, headers map[string]string) *http.Response {
	t.Helper()

	var req *http.Request
	if useJSON {
		json, err := json.Marshal(params)
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, path, bytes.NewReader(json))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		formData := url.Values{}
		for k, v := range params {
			formData.Set(k, fmt.Sprintf("%v", v))
		}
		req = httptest.NewRequest(http.MethodPost, path, strings.NewReader(formData.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	srv := newMockServer()
	srv.ServeHTTP(rec, req)

	return rec.Result()
}

// Tests that the server returns a valid JWT token when a user logs in.
func TestLogin(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"identity": model.MockApplicant.Email,
		"password": model.MockPassword,
		"role":     model.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	obj := model.SuccessResponse{}
	body, _ := io.ReadAll(res.Body)

	// Parse the response
	require.NoError(t, json.Unmarshal(body, &obj))
	require.NotEqual(t, "", obj.Token, "Response does not contain token")

	claims := model.UserClaims{}
	// Parse the embedded JWT token
	_, err := jwt.ParseWithClaims(obj.Token, &claims, mockKeyFunc)

	require.NoError(t, err)
	require.Equal(t, model.MockApplicant.Email, claims.Email)
	require.Equal(t, model.MockApplicant.Role, claims.Role)
}

// Tests that the server returns MISSING_PARAMETERS when API caller is missing required parameters.
func TestMissingParameters(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"password": model.MockPassword,
		"role":     model.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "MISSING_PARAMETERS", obj.Error)
}

// Tests that the server does not return MISSING_PARAMETERS when API caller is missing optional parameters.
func TestOptionalParameters(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"identity": model.MockApplicant.Email,
		"password": model.MockPassword,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
}

// Tests that the server returns WRONG_IDENTITY when user does not exist.
func TestLoginMissingUser(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"identity": "doesnotexist@example.com",
		"password": "password",
		"role":     1,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_IDENTITY", obj.Error)
}

// Tests that the server returns WRONG_IDENTITY when user has a different role.
func TestLoginWrongRole(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"identity": model.MockApplicant.Email,
		"password": model.MockPassword,
		"role":     model.RoleRecruiter,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_IDENTITY", obj.Error)
}

// Tests that the server returns WRONG_PASSWORD when user has wrong password.
func TestLoginWrongPassword(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/login", map[string]any{
		"identity": model.MockApplicant.Email,
		"password": "wrong",
		"role":     model.MockApplicant.Role,
	}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "WRONG_PASSWORD", obj.Error)
}

// Tests that the server returns ALREADY_LOGGED_IN when a JWT token is set.
func TestAlreadyLoggedIn(t *testing.T) {
	t.Parallel()

	testToken, _ := service.SignUserToken(model.MockApplicant, []byte("mocksecret"))

	res := testRequest(t, "/api/login", map[string]any{
		"identity": model.MockApplicant.Email,
		"password": model.MockPassword,
		"role":     model.MockApplicant.Role,
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", testToken),
	})
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "ALREADY_LOGGED_IN", obj.Error)
}

// Tests that the server returns an error conformant with shared API rules on wrong route.
func TestWrongRoute(t *testing.T) {
	t.Parallel()

	res := testRequest(t, "/api/wrong", map[string]any{}, map[string]string{})
	defer res.Body.Close()

	require.Equal(t, http.StatusNotFound, res.StatusCode)

	obj := model.ErrorResponse{}
	body, _ := io.ReadAll(res.Body)

	require.NoError(t, json.Unmarshal(body, &obj))
	require.Equal(t, "UNKNOWN", obj.Error)
	require.NotNil(t, obj.Details)
}
