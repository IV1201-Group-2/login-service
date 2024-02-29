package logging

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Echo middleware to log requests using Logrus.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			elapsed := time.Since(start).Milliseconds()

			Logger.Infof("%s %s %s in %dms",
				c.Request().Proto,
				c.Request().Method,
				c.Request().RequestURI,
				elapsed)

			return err
		}
	}
}
