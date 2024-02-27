// The package database implements functions for connecting to the database
// and querying information about users.
package database

import (
	"bufio"
	"database/sql"
	"os"
	"strings"

	// Imports Postgres driver.
	_ "github.com/lib/pq"
	// Imports SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// The microservice maintains a single connection to the database while it is active.
var connection *sql.DB

// Opens connection and pings the database.
// If the connection fails, ErrConnectionFailed is returned.
func Open(url string, maxConn int) error {
	if connection != nil {
		return ErrConnectionFailed
	}

	// Select driver based on URL
	// For Postgres: postgres://
	// For SQLite: sqlite3://
	driver := strings.Split(url, ":")[0]
	db, err := sql.Open(driver, strings.ReplaceAll(url, driver+"://", ""))
	if err != nil {
		return ErrConnectionFailed.Wrap(err)
	}

	db.SetMaxOpenConns(maxConn)
	// Always keep a single idle connection open
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(0)

	if err = db.Ping(); err != nil {
		return ErrConnectionFailed.Wrap(err)
	}

	return nil
}

// Opens connection and loads a schema file into database.
func FromFile(url string, maxConn int, schema string) error {
	if err := Open(url, maxConn); err != nil {
		return err
	}

	file, err := os.Open(schema)
	if err != nil {
		return ErrQueryFailed.Wrap(err)
	}
	defer file.Close()

	// Read schema file line by line
	// Empty lines are ignored
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) != "" {
			_, err = connection.Exec(line)
			if err != nil {
				return ErrQueryFailed.Wrap(err)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return ErrQueryFailed.Wrap(err)
	}

	return nil
}
