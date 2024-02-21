package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		//nolint:errcheck
		godotenv.Load(".env.development")
	} else {
		//nolint:errcheck
		godotenv.Load(".env")
	}

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

	db, err := database.Connect(os.Getenv("DATABASE_URL"))
	if errors.Is(err, database.ErrConnectionMockMode) {
		srv.Logger.Print("Server is in mock mode")
	} else if err != nil {
		srv.Logger.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	service.RegisterRoutes(srv, db)
	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", port)))
}
