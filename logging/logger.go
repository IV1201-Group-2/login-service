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

func init() {
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

// Log a message that occurred in the application.
func Logf(level logrus.Level, format string, args ...any) {
	logger.Logf(level, fmt.Sprintf("[Service] %s\n", format), args...)
}

// Log a message that occurred in a handler.
func Logcf(level logrus.Level, c echo.Context, format string, args ...any) {
	logger.Logf(level, fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}
