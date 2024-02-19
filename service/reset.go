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
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrTokenNotProvided})
	}

	// Check if user provided a reset token
	claims, _ := token.Claims.(*model.CustomUserClaims)
	if claims.Usage != model.TokenUsageReset {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: model.APIErrAlreadyLoggedIn})
	}

	var params resetParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: model.APIErrMissingParameters})
	}

	err := db.UpdatePassword(claims.User.ID, params.Password)
	if err != nil {
		// TODO: Handle DB connection failure gracefully
		// This should not occur under any other condition
		c.Logger().Errorf("UpdatePassword: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}

	// Create a new token valid for auth expiry period
	newToken, err := SignUserToken(claims.User, authConfig.SigningKey)
	if err != nil {
		c.Logger().Errorf("SignUserToken: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}
	return c.JSON(http.StatusOK, model.SuccessResponse{Token: newToken})
}
