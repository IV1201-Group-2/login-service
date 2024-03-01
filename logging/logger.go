// The package logging integrates Echo with Logrus, providing consistent logging for API handlers and errors.
package logging

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// TimestampFormat is a custom format for time.Format.
const TimestampFormat = "2006-01-02 15:04"

// Logger is a Logrus instance with a customized configuration.
var logger *logrus.Logger

func initLogger() {
	out := os.Stdout
	if filename, ok := os.LookupEnv("LOG_FILE"); ok {
		out, _ = os.Open(filename)
	}

	level := logrus.InfoLevel
	if levelstr, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level, _ = logrus.ParseLevel(levelstr)
	}

	logger = &logrus.Logger{
		Out:   out,
		Level: level,
		Formatter: &logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
			TimestampFormat:           TimestampFormat,
			DisableLevelTruncation:    true,
		},
	}
}

// Log a debug message that occurred in the application.
func Debugf(format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Debugf(fmt.Sprintf("[SERVICE] %s\n", format), args...)
}

// Log an informational message that occurred in the application.
func Infof(format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Infof(fmt.Sprintf("[SERVICE] %s\n", format), args...)
}

// Log a warning message that occurred in the application.
func Warnf(format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Infof(fmt.Sprintf("[SERVICE] %s\n", format), args...)
}

// Log an error that occured in the application.
func Errorf(format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Errorf(fmt.Sprintf("[SERVICE] %s\n", format), args...)
}

// Log an error that occured in the application and exit.
func Fatalf(format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Fatalf(fmt.Sprintf("[SERVICE] %s\n", format), args...)
}

// Log a debug message that occurred in a handler.
func Debugcf(c echo.Context, format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Debugf(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an informational message that occurred in a handler.
func Infocf(c echo.Context, format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Infof(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log a warning message that occurred in a handler.
func Warncf(c echo.Context, format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Infof(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an error that occurred in a handler.
func Errorcf(c echo.Context, format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Errorf(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an error that occurred in a handler and exit.
func Fatalcf(c echo.Context, format string, args ...any) {
	if logger == nil {
		initLogger()
	}
	logger.Fatalf(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}
