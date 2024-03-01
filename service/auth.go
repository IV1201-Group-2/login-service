package service

import (
	"errors"

	"github.com/IV1201-Group-2/login-service/database"
	"github.com/IV1201-Group-2/login-service/model"
	"golang.org/x/crypto/bcrypt"
)

// A value of 10 matches the cost of the default Spring BCryptPasswordEncoder.
const passwordCost = 10

// Compares a plaintext password with a hashed password stored in the database.
func ComparePassword(plaintext string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	return err == nil
}

// Encodes a password for insertion into the database.
func HashPassword(plaintext string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(plaintext), passwordCost)
	if err != nil {
		return "", ErrBcryptError.Wrap(err)
	}
	return string(result), nil
}

// Authenticate a user with the specified identity, password and optionally role.
func AuthenticateUser(repository *database.UserRepository, identity string, password string, role *model.Role) (*model.User, error) {
	// Query the database for a user with the specified username or email.
	user, err := repository.Query(identity)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return nil, ErrWrongIdentity
		}
		return nil, err
	}

	// If a role was provided, we want to make sure the user matches expectations.
	if role != nil && *role != user.Role {
		return user, ErrWrongIdentity
	}
	// Check that user has a valid password in the database
	if user.Password == "" {
		return user, ErrMissingPassword
	}
	// Check that the correct password was provided
	if !ComparePassword(password, user.Password) {
		return user, ErrWrongPassword
	}

	return user, nil
}

// Update the password of a user in the database.
func UpdatePassword(repository *database.UserRepository, token model.UserClaims, password string) error {
	// Check if user provided a reset token
	if token.Usage != model.TokenUsageReset {
		return ErrWrongUsage
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}

	return repository.UpdatePassword(token.User.ID, hashed)
}
