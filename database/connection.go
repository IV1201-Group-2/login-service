// The package database implements functions for connecting to the database
// and querying information about users.
package database

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
)

// The microservice maintains a single connection to the database while it is active.
var connection *sql.DB

// Opens a connection and pings the database.
// If the connection succeeds, a function to close the database is returned.
// If the connection fails or is already open, ErrConnectionFailed is returned.
func Open(url string, maxConn int) (func(), error) {
	if connection != nil {
		return nil, ErrConnectionFailed
	}

	// Select driver based on URL
	// For Postgres: postgres://
	// For SQLite: sqlite3://
	driver := strings.Split(url, ":")[0]
	db, err := sql.Open(driver, strings.ReplaceAll(url, driver+"://", ""))
	if err != nil {
		return nil, ErrConnectionFailed.Wrap(err)
	}

	db.SetConnMaxIdleTime(0)
	db.SetMaxOpenConns(maxConn)

	if err = db.Ping(); err != nil {
		return nil, ErrConnectionFailed.Wrap(err)
	}

	connection = db
	return func() { connection.Close() }, nil
}

// Open a database and load schema file into database.
func OpenFromFile(url string, maxConn int, schema string) (func(), error) {
	closeDB, err := Open(url, maxConn)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(schema)
	if err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) != "" {
			_, err = connection.Exec(line)
			if err != nil {
				return nil, ErrQueryFailed.Wrap(err)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, ErrQueryFailed.Wrap(err)
	}

	return closeDB, nil
}
