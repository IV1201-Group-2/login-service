package main

import (
	"os"

	"github.com/IV1201-Group-2/login-service/api"
	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/logging"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := database.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		logging.Logf(logrus.FatalLevel, "Database init error: %v", err)
	}
	defer db.Close()

	srv, err := api.NewServer(db)
	if err != nil {
		logging.Logf(logrus.FatalLevel, "Server init error: %v", err)
	}

	logging.Logf(logrus.FatalLevel, "%v", srv.Start(":"+os.Getenv("PORT")))
}
