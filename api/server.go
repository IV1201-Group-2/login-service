package api

import (
	"github.com/IV1201-Group-2/login-service/logging"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Set up an Echo server with REST API routes.
func NewServer() (*echo.Echo, error) {
	srv := echo.New()
	srv.HTTPErrorHandler = ErrorHandler
	srv.Validator = NewValidator()

	srv.Use(logging.Middleware())
	srv.Use(middleware.Recover())
	srv.Use(middleware.CORS())

	authConfig, err := NewAuthConfig()
	if err != nil {
		return nil, err
	}

	srv.Use(echojwt.WithConfig(*authConfig))
	srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("auth", authConfig)
			return next(c)
		}
	})

	srv.POST("/api/login", Login)
	srv.POST("/api/reset", PasswordReset)

	return srv, nil
}
