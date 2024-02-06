package database

import (
	"errors"

	"github.com/IV1201-Group-2/login-service/model"
)

var ErrConnectionFailed = errors.New("connection failed")
var ErrConnectionMockMode = errors.New("database is in mock mode")

var ErrUserNotFound = errors.New("user not found in db")

type Connection interface {
	QueryUser(identity string, role model.Role) (*model.User, error)
}

type MockConnection struct{}

// Attempt to connect to PSQL database
func Connect(databaseURL string) (Connection, error) {
	if databaseURL == "mock" {
		// Caller can choose to allow mock connections
		return MockConnection{}, ErrConnectionMockMode
	}

	return nil, ErrConnectionFailed
}

// Mock implementation of DB query
func (c MockConnection) QueryUser(identity string, role model.Role) (*model.User, error) {
	var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}
	for _, user := range mockAllowedUsers {
		if (user.Username == identity || user.Email == identity) && user.Role == role {
			return &user, nil
		}
	}
	return nil, ErrUserNotFound
}
