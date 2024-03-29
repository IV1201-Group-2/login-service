// The package database implements functions for connecting to the database
// and querying information about users.
package database

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IV1201-Group-2/login-service/logging"
	sq "github.com/Masterminds/squirrel"
	// Imports Postgres driver.
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Postgres uses $1, $2, etc for placeholders.
var stmtBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Opens connection and pings the database.
// If the connection fails, ErrConnectionFailed is returned.
func Open(url string) (*sql.DB, error) {
	driver := strings.Split(url, ":")[0]
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, ErrConnectionFailed.Wrap(err)
	}

	maxConn, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNECTIONS"))
	db.SetMaxOpenConns(maxConn)
	// Always keep a single idle connection open
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(0)

	if err = db.Ping(); err != nil {
		return nil, ErrConnectionFailed.Wrap(err)
	}

	// Periodic logging of database statistics
	go func() {
		for {
			logging.Logf(logrus.DebugLevel,
				"Database statistics: in use=%d idle=%d",
				db.Stats().InUse, db.Stats().Idle)
			time.Sleep(10 * time.Second)
		}
	}()

	return db, nil
}
