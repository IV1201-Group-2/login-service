package model

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Represents the routes that a user is allowed to access.
type Role int

const (
	// This user is a recruiter that can view, accept or deny applications.
	RoleRecruiter Role = iota + 1
	// This user is an applicant that can submit new applications.
	RoleApplicant
)

// Represents a user in the database.
type User struct {
	// ID of the user in the database
	ID int `json:"id"`
	// The role of this user, represents the routes they're allowed to access
	Role Role `json:"role"`

	// If the user has a username, this will be set to a non-empty string
	Username string `json:"username,omitempty"`
	// If the user has an e-mail address, this will be set to a non-empty string
	Email string `json:"email,omitempty"`

	// Bcrypt-encoded password
	Password string `json:"-"` // Omit from JSON response
}

const (
	// This is a login token.
	TokenUsageLogin = "login"
	// This is a reset token.
	TokenUsageReset = "reset"
)

// Claims that are specific to this microservice.
type CustomClaims struct {
	Usage string `json:"usage"`
}

// Custom claims that can be read by the client and other microservices.
type CustomUserClaims struct {
	CustomClaims
	User
	jwt.RegisteredClaims
}

// Encodes a password for insertion into the database.
func HashPassword(plaintext string) (string, error) {
	// Match default cost of Spring BCryptPasswordEncoder
	result, err := bcrypt.GenerateFromPassword([]byte(plaintext), 10)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Compares a plaintext password with a hashed password stored in the database.
func ComparePassword(plaintext string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	return err == nil
}
