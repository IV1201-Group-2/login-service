package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func main() {
	srv := echo.New()

	// Universal middleware for all routes
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${remote_ip}] ${protocol} ${method} ${uri} in ${latency_human} (${status})\n",
	}))
	srv.Use(middleware.Recover())
	srv.HTTPErrorHandler = customErrorHandler

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("$PORT must be set")
	}

	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", port)))
}
