package service

import (
	"time"

	"github.com/IV1201-Group-2/login-service/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Expire tokens after one hour for security
const AuthExpiryPeriod = time.Hour

// The frontend will prefer to use query strings
// A microservice or mobile app will prefer to set the Authorization header automatically
const AuthLookupMethods = "header:Authorization:Bearer,query:token"

func newClaimsFunc(c echo.Context) jwt.Claims {
	return model.UserClaims{}
}

// Use the mock signing key
var MockAuthConfig = echojwt.Config{
	NewClaimsFunc: newClaimsFunc,
	SigningKey:    model.MockJWTSigningKey,
	TokenLookup:   AuthLookupMethods,
}

func SignTokenForUser(user model.User, config *echojwt.Config) (string, error) {
	claims := model.UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AuthExpiryPeriod)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(config.SigningKey)

	if err != nil {
		return "", err
	}

	return encodedToken, nil
}
