// The package database implements functions for connecting to the database
// and querying information about users.
package database

import (
	"database/sql"
	"errors"

	"github.com/IV1201-Group-2/login-service/model"

	_ "github.com/lib/pq"
)

// Indicates that connection to the database failed.
var ErrConnectionFailed = errors.New("connection failed")

// Indicates that a mock database is being used.
// This is a warning and can be ignored if the user is informed.
var ErrConnectionMockMode = errors.New("database is in mock mode")

// Indicates that a user with the specificed identity couldn't be found.
var ErrUserNotFound = errors.New("user not found in db")

// Represents a connection to a database.
// The connection should be closed when it's no longer being used.
type Connection interface {
	// Queries the database for a user with a specific identity and role.
	QueryUser(identity string, role model.Role) (*model.User, error)
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

	err = db.Ping()
	if err != nil {
		return nil, ErrConnectionFailed
	}

	return sqlConnection{db: db}, nil
}

const userQuery = "SELECT person_id, username, email, password FROM person WHERE (username = $1 OR email = $1) AND role_id = $2"

// SQL implementation of database query.
func (c sqlConnection) QueryUser(identity string, role model.Role) (*model.User, error) {
	var name, email, password sql.NullString
	user := &model.User{Role: role}

	row := c.db.QueryRow(userQuery, identity, role)
	err := row.Scan(&user.ID, &name, &email, &password)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	if name.Valid {
		user.Username = name.String
	}
	if email.Valid {
		user.Email = email.String
	}
	if password.Valid {
		user.Password = password.String
	}

	return user, nil
}
func (c sqlConnection) Close() error {
	return c.db.Close()
}

// Mock implementation of database query.
func (c mockConnection) QueryUser(identity string, role model.Role) (*model.User, error) {
	var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}
	for _, user := range mockAllowedUsers {
		if (user.Username == identity || user.Email == identity) && user.Role == role {
			return &user, nil
		}
	}
	return nil, ErrUserNotFound
}
func (c mockConnection) Close() error {
	return nil
}
