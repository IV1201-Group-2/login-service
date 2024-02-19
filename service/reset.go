package service

import (
	"github.com/IV1201-Group-2/login-service/database"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func PasswordReset(c echo.Context, db database.Connection, authConfig *echojwt.Config) error {
	return nil
}
