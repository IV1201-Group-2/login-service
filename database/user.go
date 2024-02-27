package database

import (
	"database/sql"
	"errors"
	"os"
	"strconv"

	"github.com/IV1201-Group-2/login-service/model"
	// Imports Postgres driver.
	_ "github.com/lib/pq"
)

// Represents a connection to a database.
// The connection should be closed when it's no longer being used.
type Connection interface {
	// Queries the database for a user with a specific identity and role.
	QueryUser(identity string) (*model.User, error)
	// Updates a user password in the database.
	UpdatePassword(id int, plaintext string) error
	// Closes the database connection.
	Close() error
}

type sqlConnection struct {
	db *sql.DB
}

type mockConnection struct{}

// Attempt to connect to Postgres database.
func Connect(databaseURL string) (Connection, error) {
	if databaseURL == "mock" {
		// Caller can choose to allow mock connections.
		return mockConnection{}, ErrConnectionMockMode
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, ErrConnectionFailed
	}

	// TODO: Temporary hotfix for Heroku max connections issue
	if maxConn, ok := os.LookupEnv("DATABASE_MAX_CONNECTIONS"); ok {
		maxConn, _ := strconv.Atoi(maxConn)
		db.SetMaxOpenConns(maxConn)
	}

	err = db.Ping()
	if err != nil {
		return nil, ErrConnectionFailed
	}

	return sqlConnection{db: db}, nil
}

const userQueryStatement = "SELECT person_id, username, email, password, role_id FROM person WHERE (username = $1 OR email = $1)"

// SQL implementation of database query.
func (c sqlConnection) QueryUser(identity string) (*model.User, error) {
	var name, email, password sql.NullString
	user := &model.User{}

	row := c.db.QueryRow(userQueryStatement, identity)
	err := row.Scan(&user.ID, &name, &email, &password, &user.Role)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	user.Username = name.String
	user.Email = email.String
	user.Password = password.String

	return user, nil
}

const updatePasswordStatement = "UPDATE person SET password = $2 WHERE person_id = $1"

// SQL implementation of database password update.
func (c sqlConnection) UpdatePassword(id int, plaintext string) error {
	hashedPassword, err := model.HashPassword(plaintext)
	if err != nil {
		return err
	}

	result, err := c.db.Exec(updatePasswordStatement, id, hashedPassword)

	// Error executing statement
	if err != nil {
		return err
	}
	// Error finding user
	if rows, err := result.RowsAffected(); err != nil || rows != 1 {
		return ErrUserNotFound
	}

	return nil
}

func (c sqlConnection) Close() error {
	return c.db.Close()
}

// Mock implementation of database query.
func (c mockConnection) QueryUser(identity string) (*model.User, error) {
	var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}
	for _, user := range mockAllowedUsers {
		if user.Username == identity || user.Email == identity {
			return &user, nil
		}
	}
	return nil, ErrUserNotFound
}

// Mock implementation of database password update. Not supported.
func (c mockConnection) UpdatePassword(_ int, _ string) error {
	return ErrUserNotFound
}

func (c mockConnection) Close() error {
	return nil
}
