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

// Custom claims that can be read by the client and other microservices.
type UserClaims struct {
	User
	jwt.RegisteredClaims
}

// Compares a plaintext password with a hashed password stored in the database.
func ComparePassword(plaintext string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	return err == nil
}
