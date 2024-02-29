package service

import (
	"time"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/golang-jwt/jwt/v5"
)

// Expire tokens after one hour for security.
const TokenExpiryPeriod = time.Hour

// Expire reset tokens after ten minutes.
const TokenResetExpiryPeriod = time.Minute * 10

func signToken(claims jwt.Claims, signingKey any) (string, time.Time, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiry, _ := claims.GetExpirationTime()

	encodedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", time.Now(), ErrJWTError.Wrap(err)
	}

	return encodedToken, expiry.Time, nil
}

// Signs a token for the specified user with the specified signing key.
// This function returns the encoded token in plaintext or an error if signing failed.
func SignUserToken(user model.User, signingKey any) (string, time.Time, error) {
	claims := model.UserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageLogin,
		},
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiryPeriod)),
		},
	}
	return signToken(claims, signingKey)
}

// Signs a reset token for the specified user with the specified signing key.
// The reset token should be sent to the user through a secure channel (such as email)
// since it grants temporary access to an account without a password.
// This function returns the encoded token in plaintext or an error if signing failed.
func SignResetToken(user model.User, signingKey any) (string, time.Time, error) {
	claims := model.UserClaims{
		CustomClaims: model.CustomClaims{
			Usage: model.TokenUsageReset,
		},
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenResetExpiryPeriod)),
		},
	}
	return signToken(claims, signingKey)
}
