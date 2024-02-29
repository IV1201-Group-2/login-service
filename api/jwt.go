package api

import (
	"errors"
	"os"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func errorHandlerFunc(_ echo.Context, err error) error {
	// Allow requests without a token set
	if errors.Is(err, echojwt.ErrJWTMissing) {
		return nil
	}
	if errors.Is(err, echojwt.ErrJWTInvalid) {
		return ErrTokenInvalid.Wrap(err)
	}

	return err
}

func newClaimsFunc(_ echo.Context) jwt.Claims {
	return &model.UserClaims{}
}

var authConfigTemplate = echojwt.Config{
	ErrorHandler:           errorHandlerFunc,
	ContinueOnIgnoredError: true,
	NewClaimsFunc:          newClaimsFunc,
}

// ErrNoSecret indicates that the JWT_SECRET environment variable is not set.
var ErrNoSecret = errors.New("$JWT_SECRET must be set")

// NewAuthConfig creates a new echojwt config from JWT_SECRET.
func NewAuthConfig() (*echojwt.Config, error) {
	if secret, ok := os.LookupEnv("JWT_SECRET"); ok {
		config := authConfigTemplate
		config.SigningKey = []byte(secret)
		return &config, nil
	}
	return nil, ErrNoSecret
}
