package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	srv := echo.New()

	srv.HTTPErrorHandler = service.ErrorHandler
	srv.Validator = service.NewValidator()

	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${remote_ip}] ${protocol} ${method} ${uri} in ${latency_human} (${status})\n",
	}))
	srv.Logger.SetHeader("")

	srv.Use(middleware.Recover())
	srv.Use(middleware.CORS())

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		srv.Logger.Fatal("$PORT must be set")
	}

	maxConn, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNECTIONS"))
	closeDB, err := database.Open(os.Getenv("DATABASE_URL"), maxConn)
	if err != nil {
		srv.Logger.Fatalf("Database error: %v", err)
	}
	defer closeDB()

	service.RegisterRoutes(srv)
	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", port)))
}
