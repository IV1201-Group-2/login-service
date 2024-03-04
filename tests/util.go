package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Database is a single database connection that is maintained for all tests.
var Database *sql.DB

// Set up an appropriate environment for testing.
// If this function succeeds, it returns a cleanup function.
func SetupEnvironment() (func() error, error) {
	os.Setenv("DATABASE_MAX_CONNECTIONS", "8")
	os.Setenv("JWT_SECRET", MockSecret)

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
	Database, err = database.Open(connStr)
	if err != nil {
		return nil, err
	}

	return func() error {
		if Database != nil {
			Database.Close()
		}
		err := pgContainer.Terminate(context.Background())
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// Sends a request to an existing server and returns the response.
func CustomRequest(t *testing.T, srv *echo.Echo, path string, params map[string]any, headers map[string]string) *http.Response {
	t.Helper()

	json, err := json.Marshal(params)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec.Result()
}

// Sends a request to a mock server and returns the response.
func Request(t *testing.T, path string, params map[string]any, headers map[string]string) *http.Response {
	t.Helper()

	srv, _ := api.NewServer(Database)
	defer srv.Close()

	return CustomRequest(t, srv, path, params, headers)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// Generate a random string of a fixed length.
func RandomStr(length int) string {
	src := rand.New(rand.NewSource(time.Now().UnixMilli())) // #nosec G404
	str := make([]rune, length)
	for i := 0; i < length; i++ {
		str[i] = letters[src.Intn(len(letters))]
	}
	return string(str)
}
