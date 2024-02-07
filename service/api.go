// The package service contains the implementation of the microservice.
// This includes handling the login process and signing JWT tokens.
package service

import (
	"errors"
	"log"
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

		// Something went wrong, DB down?
		c.Logger().Errorf("/api/login QueryUser: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}

	// If the caller specified a role, we want to check if the user matches expectations
	if params.Role > 0 && params.Role != user.Role {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrWrongIdentity})
	}

	// TODO: User should be sent an email
	if user.Password == "" {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: model.APIErrMissingPassword})
	}
	// Check that password matches
	if !model.ComparePassword(params.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: model.APIErrWrongPassword})
	}

	// Create a new token valid for auth expiry period
	token, err := SignTokenForUser(*user, authConfig.SigningKey)
	if err != nil {
		// Something went wrong when signing the token
		c.Logger().Errorf("/api/login SignTokenForUser: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: model.APIErrUnknown})
	}

	return c.JSON(http.StatusOK, model.SuccessResponse{Token: token})
}

// Register all REST API routes.
func RegisterRoutes(srv *echo.Echo, db database.Connection) {
	authConfig, err := NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv.Use(echojwt.WithConfig(*authConfig))
	srv.POST("/api/login", func(c echo.Context) error {
		return Login(c, db, authConfig)
	})
}
