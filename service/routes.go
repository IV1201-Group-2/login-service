package service

import (
	"log"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Register all REST API routes.
func RegisterRoutes(srv *echo.Echo) {
	authConfig, err := NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv.Use(echojwt.WithConfig(*authConfig))
	srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("authConfig", authConfig)
			return next(c)
		}
	})

	srv.POST("/api/login", Login)
	srv.POST("/api/reset", PasswordReset)
}
