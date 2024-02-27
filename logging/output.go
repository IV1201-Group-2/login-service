// The package logging integrates Echo with Logrus, providing consistent logging for API handlers and errors.
package logging

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Format for time.Format
const TimestampFormat = "2006-01-02 15:04"

// Customized configuration for Logrus.
var Logger = logrus.Logger{
	Out: os.Stdout,
	Formatter: &logrus.TextFormatter{
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           TimestampFormat,
		DisableLevelTruncation:    true,
	},
	Level: logrus.InfoLevel,
}

// Log an informational message that occurred in a handler.
func Infof(c echo.Context, format string, args ...any) {
	Logger.Logf(logrus.InfoLevel, fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an error that occurred in a handler.
func Errorf(c echo.Context, format string, args ...any) {
	Logger.Logf(logrus.ErrorLevel, fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an error that occurred in a handler and panic.
func Panicf(c echo.Context, format string, args ...any) {
	Logger.Logf(logrus.PanicLevel, fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}
