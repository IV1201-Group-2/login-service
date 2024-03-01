package database_test

import (
	"os"
	"testing"

	"github.com/IV1201-Group-2/login-service/logging"
	"github.com/IV1201-Group-2/login-service/tests"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	cleanup, err := tests.SetupEnvironment()
	if err != nil {
		logging.Logf(logrus.ErrorLevel, "Failed to set up test environment: %v", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	err = cleanup()
	if err != nil {
		logging.Logf(logrus.ErrorLevel, "Failed to tear down test environment: %v", err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}
