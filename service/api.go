package service

import (
	"errors"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type loginParams struct {
	Identity string     `form:"identity" validate:"required"`
	Password string     `form:"password" validate:"required"`
	Role     model.Role `form:"role" validate:"required"`
}

// Login route handler
func Login(c echo.Context, db database.Connection) error {
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

	if user, err := db.QueryUser(params.Identity, params.Role); err != nil {
		if err == database.ErrUserNotFound {
			return c.JSON(401, model.ErrorResponse{Error: model.APIErrWrongIdentity})
		} else {
			// Something went wrong, DB down?
			c.Logger().Errorf("/login: %v", err)
			return c.JSON(500, model.ErrorResponse{Error: model.APIErrUnknown})
		}
	} else {
		// Check that password matches
		if user.Password != params.Password {
			return c.JSON(401, model.ErrorResponse{Error: model.APIErrWrongPassword})
		}

		// Create a new token valid for auth expiry period
		if token, err := SignTokenForUser(*user, &AuthConfig); err == nil {
			return c.JSON(200, model.SuccessResponse{Token: token})
		} else {
			// Something went wrong when signing the token
			c.Logger().Errorf("/login: %v", err)
			return c.JSON(500, model.ErrorResponse{Error: model.APIErrUnknown})
		}
	}
}

func RegisterRoutes(srv *echo.Echo, db database.Connection) {
	srv.Use(echojwt.WithConfig(AuthConfig))
	srv.POST("/api/login", func(c echo.Context) error {
		return Login(c, db)
	})
}
