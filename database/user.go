package database

import (
	"errors"

	"github.com/IV1201-Group-2/login-service/model"
)

var ErrUserNotFound = errors.New("user not found in db")

// Mock implementation of user DB query
func MockQueryUser(identity string, role model.Role) (*model.User, error) {
	var mockAllowedUsers = []model.User{model.MockApplicant, model.MockRecruiter}

	for _, user := range mockAllowedUsers {
		if user.Role == role && (user.Username == identity || user.Email == identity) {
			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}
