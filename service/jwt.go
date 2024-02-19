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

// Expire reset tokens after ten minutes.
const AuthResetExpiryPeriod = time.Minute * 10

func errorHandlerFunc(_ echo.Context, err error) error {
	// Allow requests without a token set
	if errors.Is(err, echojwt.ErrJWTMissing) {
		return nil
	}

	return err
}

func newClaimsFunc(_ echo.Context) jwt.Claims {
	return &model.CustomUserClaims{}
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

func signToken(claims jwt.Claims, signingKey interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return encodedToken, nil
}

// Signs a token for the specified user with the specified signing key.
// This function returns the encoded token in plaintext or an error if signing failed.
func SignUserToken(user model.User, signingKey interface{}) (string, error) {
	claims := model.CustomUserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageLogin,
		},
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AuthExpiryPeriod)),
		},
	}
	return signToken(claims, signingKey)
}

// Signs a reset token for the specified user with the specified signing key.
// The reset token should be sent to the user through a secure channel (such as email)
// since it grants temporary access to an account without a password.
// This function returns the encoded token in plaintext or an error if signing failed.
func SignResetToken(user model.User, signingKey interface{}) (string, error) {
	claims := model.CustomUserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageReset,
		},
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AuthResetExpiryPeriod)),
		},
	}
	return signToken(claims, signingKey)
}
