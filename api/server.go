package api

import (
	"database/sql"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer creates a new Echo server instance for the login REST API.
func NewServer(db *sql.DB) (*echo.Echo, error) {
	srv := echo.New()
	srv.HTTPErrorHandler = ErrorHandler
	srv.Validator = NewValidator()

	srv.Use(logging.Middleware())
	srv.Use(middleware.Recover())
	srv.Use(middleware.CORS())

	userRepository := database.NewUserRepository(db)

	authConfig, err := NewAuthConfig()
	if err != nil {
		return nil, err
	}
	srv.Use(echojwt.WithConfig(*authConfig))

	srv.POST("/api/login", func(c echo.Context) error {
		return Login(c, userRepository, authConfig)
	})
	srv.POST("/api/reset", func(c echo.Context) error {
		return PasswordReset(c, userRepository, authConfig)
	})

	return srv, nil
}
