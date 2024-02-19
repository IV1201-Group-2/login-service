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
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: model.APIErrAlreadyLoggedIn})
	}

	var params loginParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: model.APIErrMissingParameters})
	}

	user, err := db.QueryUser(params.Identity)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrWrongIdentity})
		}
		// TODO: Handle DB connection failure gracefully
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}

	// If the caller specified a role, we want to check if the user matches expectations
	if params.Role > 0 && params.Role != user.Role {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrWrongIdentity})
	}

	// Check that user has a valid password in the database
	if user.Password == "" {
		// Create a new reset token allowing the user to reset their password
		token, err := SignResetToken(*user, authConfig.SigningKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
		}
		response := model.ErrorResponse{
			Error: model.APIErrMissingPassword,
			Details: &model.ErrorDetails{
				ResetToken: token,
			},
		}
		return c.JSON(http.StatusNotFound, response)
	}

	// Check that password matches
	if !model.ComparePassword(params.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrWrongPassword})
	}

	// Create a new token valid for auth expiry period
	token, err := SignUserToken(*user, authConfig.SigningKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}
	return c.JSON(http.StatusOK, model.SuccessResponse{Token: token})
}
