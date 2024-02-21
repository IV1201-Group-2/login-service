package service

import (
	"errors"
	"net/http"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type resetParams struct {
	Password string `form:"password" json:"password" validate:"required"`
}

// Password reset route handler.
func PasswordReset(c echo.Context, db database.Connection, authConfig *echojwt.Config) error {
	// Check if user provided a token
	// TODO: Verify that Echo checks expiry period
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		LogErrorf(c, "Unauthorized attempt: user has no reset token")
		return model.ErrTokenNotProvided
	}

	// Check if user provided a reset token
	claims, _ := token.Claims.(*model.CustomUserClaims)
	if claims.Usage != model.TokenUsageReset {
		LogErrorf(c, "Unauthorized attempt: user has no reset token")
		return model.ErrAlreadyLoggedIn
	}

	var params resetParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return model.ErrMissingParameters
	}

	err := db.UpdatePassword(claims.User.ID, params.Password)
	if err != nil {
		return err
	}

	if claims.User.Username != "" {
		LogErrorf(c, "User '%s' has reset password", claims.User.Username)
	} else if claims.User.Email != "" {
		LogErrorf(c, "User '%s' has reset password", claims.User.Email)
	}

	// Create a new token valid for auth expiry period
	newToken, expiry, err := SignUserToken(claims.User, authConfig.SigningKey)
	if err != nil {
		return err
	}
	LogErrorf(c, "Login successful: token expires at %s", expiry.Format("2006-01-02 15:04"))

	return c.JSON(http.StatusOK, model.UserTokenResponse{Token: newToken})
}
