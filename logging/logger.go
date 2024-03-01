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
var Logger = logrus.Logger{
	Out:   os.Stdout,
	Level: logrus.InfoLevel,
	Formatter: &logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           TimestampFormat,
		DisableLevelTruncation:    true,
	},
}

// Log an informational message that occurred in a handler.
func Infof(c echo.Context, format string, args ...any) {
	Logger.Infof(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log a warning message that occurred in a handler.
func Warnf(c echo.Context, format string, args ...any) {
	Logger.Infof(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}

// Log an error that occurred in a handler.
func Errorf(c echo.Context, format string, args ...any) {
	Logger.Errorf(fmt.Sprintf("[%s] %s\n", c.RealIP(), format), args...)
}
