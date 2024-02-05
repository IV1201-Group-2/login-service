package service

import (
	"errors"

	"github.com/IV1201-Group-2/login-service/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}

// TODO: Discuss query params for POST
type loginParams struct {
	Identity string     `json:"identity" validate:"required"`
	Password string     `json:"password" validate:"required"`
	Role     model.Role `json:"role" validate:"required"`
}

// Mock implementation of login API
func MockLogin(c echo.Context) error {
	// Check if user incorrectly provided a JWT token
	_, ok := c.Get("user").(*jwt.Token)
	if ok {
		return c.JSON(400, model.ErrorResponse{Error: model.APIErrAlreadyLoggedIn})
	}

	var params loginParams
	// Check that all parameters are present
	if err := errors.Join(c.Bind(&params), c.Validate(&params)); err != nil {
		return c.JSON(400, model.ErrorResponse{Error: model.APIErrMissingParameters})
	}

	// Check if user is in list of mock users
	for _, user := range mockAllowedUsers {
		if user.Role == params.Role && (user.Username == params.Identity || user.Email == params.Identity) {
			// Check that password matches
			if user.Password != params.Password {
				return c.JSON(401, model.ErrorResponse{Error: model.APIErrWrongPassword})
			}

			// Create a new token valid for auth expiry period
			if token, err := SignTokenForUser(user, &MockAuthConfig); err == nil {
				return c.JSON(200, model.LoginSuccessResponse{Token: token})
			} else {
				// Something went wrong when signing the token
				c.Logger().Errorf("/login: %v", err)
				return c.JSON(500, model.ErrorResponse{Error: model.APIErrUnknown})
			}
		}
	}

	return c.JSON(401, model.ErrorResponse{Error: model.APIErrWrongIdentity})
}

func RegisterMockRoutes(srv *echo.Echo) {
	srv.StdLogger.Println("Server is in mock mode")
	srv.Use(echojwt.WithConfig(MockAuthConfig))
	srv.POST("/login", MockLogin)
}
