package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/IV1201-Group-2/login-service/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"
)

// Custom error handler conformant with shared API rules
// https://echo.labstack.com/docs/error-handling
func customErrorHandler(err error, c echo.Context) {
	var details *model.ErrorDetails
	code := http.StatusInternalServerError

	c.Logger().Error(err)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		details = &model.ErrorDetails{Message: fmt.Sprintf("%v", httpErr.Message)}
		code = httpErr.Code
	}
	response := model.ErrorResponse{
		Error:   model.APIErrUnknown,
		Details: details,
	}

	if err := c.JSON(code, response); err != nil {
		c.Logger().Error(err)
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, fmt.Sprintf("validation failed: %s", err.Error()))
	}

	return nil
}

func main() {
	if os.Getenv("APP_ENV") == "development" {
		godotenv.Load(".env.development")
	} else {
		godotenv.Load(".env")
	}

	srv := echo.New()

	// Universal middleware for all routes
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${remote_ip}] ${protocol} ${method} ${uri} in ${latency_human} (${status})\n",
	}))
	srv.Use(middleware.Recover())

	srv.HTTPErrorHandler = customErrorHandler
	srv.Validator = &CustomValidator{validator: validator.New()}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("$PORT must be set")
	}

	service.RegisterMockRoutes(srv)
	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", port)))
}
