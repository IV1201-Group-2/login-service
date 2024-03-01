package main

import (
	"os"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
)

func main() {
	db, err := database.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		logging.Fatalf("Database init error: %v", err)
	}
	defer db.Close()

	srv, err := api.NewServer(db)
	if err != nil {
		logging.Fatalf("Server init error: %v", err)
	}

	logging.Fatalf("%v", srv.Start(":"+os.Getenv("PORT")))
}
