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

type loginParams struct {
	Identity string     `form:"identity" json:"identity" validate:"required"`
	Password string     `form:"password" json:"password" validate:"required"`
	Role     model.Role `form:"role"     json:"role"     validate:"omitempty"`
}

// Login route handler.
func Login(c echo.Context, db database.Connection, authConfig *echojwt.Config) error {
	// Check if user incorrectly provided a JWT token
	_, ok := c.Get("user").(*jwt.Token)
	if ok {
		return model.ErrAlreadyLoggedIn
	}

	var params loginParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return model.ErrMissingParameters
	}

	user, err := db.QueryUser(params.Identity)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			LogErrorf(c, "Unauthorized attempt: user '%s' not found", params.Identity)
			return model.ErrWrongIdentity
		}
		return err
	}

	// If the caller specified a role, we want to check if the user matches expectations
	if params.Role > 0 && params.Role != user.Role {
		return model.ErrWrongIdentity
	}

	// Check that user has a valid password in the database
	if user.Password == "" {
		// Create a new reset token allowing the user to reset their password
		token, expiry, err := SignResetToken(*user, authConfig.SigningKey)
		if err != nil {
			return err
		}
		LogErrorf(c, "Login failed: user has no password in db")
		LogErrorf(c, "Handed out reset token that expires at %s", expiry.Format("2006-01-02 15:04"))
		return model.ErrMissingPassword.WithDetails(model.ResetTokenResponse{ResetToken: token})
	}

	// Check that password matches
	if !model.ComparePassword(params.Password, user.Password) {
		LogErrorf(c, "Unauthorized attempt: wrong password for user '%s'", params.Identity)
		return model.ErrWrongPassword
	}

	// Create a new token valid for auth expiry period
	token, expiry, err := SignUserToken(*user, authConfig.SigningKey)
	if err != nil {
		return err
	}
	LogErrorf(c, "Login successful: token expires at %s", expiry.Format("2006-01-02 15:04"))

	return c.JSON(http.StatusOK, model.UserTokenResponse{Token: token})
}
