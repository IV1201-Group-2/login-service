package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
)

func main() {
	maxConn, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNECTIONS"))
	err := database.Open(os.Getenv("DATABASE_URL"), maxConn)
	if err != nil {
		logging.Logger.Fatalf("Database error: %v", err)
	}

	srv, err := api.NewServer()
	if err != nil {
		logging.Logger.Fatalf("Server error: %v", err)
	}

	logging.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", os.Getenv("PORT"))))
}
