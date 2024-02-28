// The package database implements functions for connecting to the database
// and querying information about users.
package database

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"

	// Imports Postgres driver.
	_ "github.com/lib/pq"
)

// Postgres uses $1, $2, etc for placeholders
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

	return db, nil
}
