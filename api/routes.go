// The package api is responsible for translating requests and responses for presentation to the user.
package api

import (
	"errors"
	"net/http"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type loginParams struct {
	Identity string      `form:"identity" json:"identity" validate:"required"`
	Password string      `form:"password" json:"password" validate:"required"`
	Role     *model.Role `form:"role"     json:"role"     validate:"omitempty"`
}

// Login route handler.
func Login(c echo.Context, userRepository *database.UserRepository, auth *echojwt.Config) error {
	// Check if user incorrectly provided a JWT token
	_, ok := c.Get("user").(*jwt.Token)
	if ok {
		return ErrAlreadyLoggedIn
	}

	var params loginParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return ErrMissingParameters
	}

	user, err := service.AuthenticateUser(userRepository, params.Identity, params.Password, params.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMissingPassword):
			// Create a new reset token allowing the user to reset their password
			token, expiry, err := service.SignResetToken(*user, auth.SigningKey)
			if err != nil {
				return err
			}
			logging.Logcf(logrus.WarnLevel, c, "Login failed: user has no password in db")
			logging.Logcf(logrus.InfoLevel, c, "Handed out reset token that expires at %s", expiry.Format(logging.TimestampFormat))
			return ErrMissingPassword.WithDetails(model.ResetTokenResponse{Token: token})
		case errors.Is(err, service.ErrWrongIdentity):
			logging.Logcf(logrus.WarnLevel, c, "Unauthorized attempt: user '%s' not found", params.Identity)
			return ErrWrongIdentity
		case errors.Is(err, service.ErrWrongPassword):
			logging.Logcf(logrus.WarnLevel, c, "Unauthorized attempt: wrong password for user '%s'", params.Identity)
			return ErrWrongIdentity
		}

		return err
	}

	// Create a new token valid for the auth expiry period
	token, expiry, err := service.SignUserToken(*user, auth.SigningKey)
	if err != nil {
		return err
	}
	logging.Logcf(logrus.InfoLevel, c, "Login successful: token expires at %s", expiry.Format(logging.TimestampFormat))

	return c.JSON(http.StatusOK, model.LoginTokenResponse{Token: token})
}

type resetParams struct {
	Password string `form:"password" json:"password" validate:"required"`
}

// Password reset route handler.
func PasswordReset(c echo.Context, userRepository *database.UserRepository, auth *echojwt.Config) error {
	// Check if user provided a token
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		logging.Logcf(logrus.WarnLevel, c, "Unauthorized attempt: user has no reset token")
		return ErrTokenNotProvided
	}

	var params resetParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return ErrMissingParameters
	}

	claims, _ := token.Claims.(*model.UserClaims)
	err := service.UpdatePassword(userRepository, *claims, params.Password)
	if errors.Is(err, service.ErrWrongUsage) {
		return ErrTokenInvalid
	} else if err != nil {
		return err
	}

	if claims.User.Username != "" {
		logging.Logcf(logrus.InfoLevel, c, "User '%s' has reset password", claims.User.Username)
	} else if claims.User.Email != "" {
		logging.Logcf(logrus.InfoLevel, c, "User '%s' has reset password", claims.User.Email)
	}

	// Create a new token valid for the auth expiry period
	newToken, expiry, err := service.SignUserToken(claims.User, auth.SigningKey)
	if err != nil {
		return err
	}
	logging.Logcf(logrus.InfoLevel, c, "Login successful: token expires at %s", expiry.Format(logging.TimestampFormat))

	return c.JSON(http.StatusOK, model.LoginTokenResponse{Token: newToken})
}
