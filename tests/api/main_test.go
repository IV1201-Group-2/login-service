package api_test

import (
	"os"
	"testing"

	"github.com/IV1201-Group-2/login-service/logging"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	cleanup, err := tests.SetupEnvironment()
	if err != nil {
		logging.Logger.Errorf("Failed to set up test environment: %v", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	err = cleanup()
	if err != nil {
		logging.Logger.Errorf("Failed to tear down test environment: %v", err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func mockKeyFunc(_ *jwt.Token) (any, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}
