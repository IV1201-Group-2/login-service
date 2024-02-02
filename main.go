package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	srv := echo.New()

	// Universal middleware for all routes
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${remote_ip}] ${protocol} ${method} ${uri} in ${latency_human} (${status})\n",
	}))
	srv.Use(middleware.Recover())

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("$PORT must be set")
	}

	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", port)))
}
