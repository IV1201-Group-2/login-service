package service

import (
	"log"

	"github.com/IV1201-Group-2/login-service/database"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Register all REST API routes.
func RegisterRoutes(srv *echo.Echo, db database.Connection) {
	authConfig, err := NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}
	srv.Use(echojwt.WithConfig(*authConfig))

	srv.POST("/api/login", func(c echo.Context) error {
		return Login(c, db, authConfig)
	})
	srv.POST("/api/reset", func(c echo.Context) error {
		return PasswordReset(c, db, authConfig)
	})
}
