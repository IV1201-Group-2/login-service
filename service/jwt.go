package service

import (
	"errors"
	"os"
	"time"

	"github.com/IV1201-Group-2/login-service/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Expire tokens after one hour for security.
const AuthExpiryPeriod = time.Hour

func errorHandlerFunc(c echo.Context, err error) error {
	// Allow requests without a token set
	if errors.Is(err, echojwt.ErrJWTMissing) {
		return nil
	}

	return err
}

func newClaimsFunc(c echo.Context) jwt.Claims {
	return &model.UserClaims{}
}

var authConfigTemplate = echojwt.Config{
	ErrorHandler:           errorHandlerFunc,
	ContinueOnIgnoredError: true,
	NewClaimsFunc:          newClaimsFunc,
}

// Indicates that the JWT_SECRET environment variable is not set.
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

// Signs a token for the specified user with the specified authentication config.
// This function returns the encoded token in plaintext or an error if signing failed.
func SignTokenForUser(user model.User, signingKey interface{}) (string, error) {
	claims := model.UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AuthExpiryPeriod)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return encodedToken, nil
}
