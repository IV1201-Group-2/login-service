package database

import (
	"database/sql"
	"errors"

	"github.com/IV1201-Group-2/login-service/model"

	_ "github.com/lib/pq"
)

var ErrConnectionFailed = errors.New("connection failed")
var ErrConnectionMockMode = errors.New("database is in mock mode")

var ErrUserNotFound = errors.New("user not found in db")

type Connection interface {
	// Queries the database for a user with a specific identity and role
	QueryUser(identity string, role model.Role) (*model.User, error)
	// Closes the database connection
	Close() error
}

type SQLConnection struct {
	db *sql.DB
}

type MockConnection struct{}

// Attempt to connect to Postgres database
func Connect(databaseURL string) (Connection, error) {
	if databaseURL == "mock" {
		// Caller can choose to allow mock connections
		return MockConnection{}, ErrConnectionMockMode
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, ErrConnectionFailed
	}

	err = db.Ping()
	if err != nil {
		return nil, ErrConnectionFailed
	}

	return SQLConnection{db: db}, nil
}

const userQuery = "SELECT person_id, username, email, password FROM person WHERE (username = $1 OR email = $1) AND role_id = $2"

// SQL implementation of database query
func (c SQLConnection) QueryUser(identity string, role model.Role) (*model.User, error) {
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
func (c SQLConnection) Close() error {
	return c.db.Close()
}

// Mock implementation of database query
func (c MockConnection) QueryUser(identity string, role model.Role) (*model.User, error) {
	var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}
	for _, user := range mockAllowedUsers {
		if (user.Username == identity || user.Email == identity) && user.Role == role {
			return &user, nil
		}
	}
	return nil, ErrUserNotFound
}
func (c MockConnection) Close() error {
	return nil
}
