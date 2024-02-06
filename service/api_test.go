package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// newMockServer creates a new Echo server that uses the mock dataset
func newMockServer() *echo.Echo {
	srv := echo.New()
	db, _ := database.Connect("mock")

	srv.HTTPErrorHandler = ErrorHandler
	srv.Validator = NewValidator()

	// TODO: Maybe better to pass a new auth config instead...
	os.Setenv("JWT_SECRET", "mocksecret")
	RegisterRoutes(srv, db)

	return srv
}

func mockKeyFunc(t *jwt.Token) (interface{}, error) {
	return []byte("mocksecret"), nil
}

// Sends a request to mock server and returns response
func testRequest(path string, params map[string]string, headers map[string]string) *http.Response {
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
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
	res := testRequest("/api/login", map[string]string{
		"identity": model.MockApplicant.Email,
		"password": model.MockApplicant.Password,
		"role":     strconv.Itoa(int(model.MockApplicant.Role)),
	}, map[string]string{})

	if assert.Equal(t, 200, res.StatusCode) {
		obj := model.SuccessResponse{}
		body, _ := io.ReadAll(res.Body)

		// Parse the response
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.NotEqual(t, "", obj.Token, "Response does not contain token")

		claims := model.UserClaims{}
		// Parse the embedded JWT token
		_, err := jwt.ParseWithClaims(obj.Token, &claims, mockKeyFunc)

		assert.NoError(t, err)
		assert.Equal(t, model.MockApplicant.Email, claims.Email)
		assert.Equal(t, model.MockApplicant.Role, claims.Role)
	}
}

// Tests that the server returns MISSING_PARAMETERS when API caller is missing parameters.
func TestMissingParameters(t *testing.T) {
	res := testRequest("/api/login", map[string]string{
		"password": model.MockApplicant.Password,
		"role":     strconv.Itoa(int(model.MockApplicant.Role)),
	}, map[string]string{})

	if assert.Equal(t, 400, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "MISSING_PARAMETERS", obj.Error)
	}
}

// Tests that the server returns WRONG_IDENTITY when user does not exist.
func TestLoginMissingUser(t *testing.T) {
	res := testRequest("/api/login", map[string]string{
		"identity": "doesnotexist@example.com",
		"password": "password",
		"role":     "1",
	}, map[string]string{})

	if assert.Equal(t, 401, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "WRONG_IDENTITY", obj.Error)
	}
}

// Tests that the server returns WRONG_IDENTITY when user has a different role.
func TestLoginWrongRole(t *testing.T) {
	res := testRequest("/api/login", map[string]string{
		"identity": model.MockApplicant.Email,
		"password": model.MockApplicant.Password,
		"role":     strconv.Itoa(int(model.RoleRecruiter)),
	}, map[string]string{})

	if assert.Equal(t, 401, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "WRONG_IDENTITY", obj.Error)
	}
}

// Tests that the server returns WRONG_PASSWORD when user has wrong password.
func TestLoginWrongPassword(t *testing.T) {
	res := testRequest("/api/login", map[string]string{
		"identity": model.MockApplicant.Email,
		"password": "wrong",
		"role":     strconv.Itoa(int(model.MockApplicant.Role)),
	}, map[string]string{})

	if assert.Equal(t, 401, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "WRONG_PASSWORD", obj.Error)
	}
}

// Tests that the server returns ALREADY_LOGGED_IN when a JWT token is set.
func TestAlreadyLoggedIn(t *testing.T) {
	testToken, _ := SignTokenForUser(model.MockApplicant, []byte("mocksecret"))

	res := testRequest("/api/login", map[string]string{
		"identity": model.MockApplicant.Email,
		"password": model.MockApplicant.Password,
		"role":     strconv.Itoa(int(model.MockApplicant.Role)),
	}, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", testToken),
	})

	if assert.Equal(t, 400, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "ALREADY_LOGGED_IN", obj.Error)
	}
}

// Tests that the server returns an error conformant with shared API rules on wrong route.
func TestWrongRoute(t *testing.T) {
	res := testRequest("/api/wrong", map[string]string{}, map[string]string{})

	if assert.Equal(t, 404, res.StatusCode) {
		obj := model.ErrorResponse{}
		body, _ := io.ReadAll(res.Body)
		assert.NoError(t, json.Unmarshal(body, &obj))
		assert.Equal(t, "UNKNOWN", obj.Error)
		assert.NotNil(t, obj.Details)
	}
}
