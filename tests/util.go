package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// The tests maintain a single connection to the database.
var testDB *sql.DB

// If this is set to true, JSON blobs are sent to the server when a test request is executed.
// If this is set to false, form data is sent to the server when a test request is executed.
var useJSON = os.Getenv("TEST_JSON") == "1"

// Set up an appropriate environment for testing.
// If this function succeeds, it returns a cleanup function.
func SetupEnvironment() (func() error, error) {
	if testDB != nil {
		return nil, errors.New("environment already set up")
	}

	// Set up a Postgres container with our test schema
	// https://testcontainers.com/guides/getting-started-with-testcontainers-for-go
	pgContainer, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithInitScripts("../schema.sql"),
		postgres.WithDatabase("public"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		return nil, err
	}
	testDB, err = database.Open(connStr)
	if err != nil {
		return nil, err
	}

	os.Setenv("JWT_SECRET", MockSecret)
	logging.Logger.SetOutput(ioutil.Discard)

	return func() error {
		if testDB != nil {
			testDB.Close()
		}
		err := pgContainer.Terminate(context.Background())
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// Sends a request to a mock server and returns the response.
func Request(t *testing.T, path string, params map[string]any, headers map[string]string) *http.Response {
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

	srv, _ := api.NewServer(testDB)
	srv.ServeHTTP(rec, req)

	return rec.Result()
}
