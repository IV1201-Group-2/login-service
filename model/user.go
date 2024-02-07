package model

import "github.com/golang-jwt/jwt/v5"

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
	ID   int  `json:"id"`
	Role Role `json:"role"`

	Username string `json:"username"`
	Email    string `json:"email"`

	// Omit from JSON response
	Password string `json:"-"`
}

// Custom claims that can be read by the client and other microservices.
type UserClaims struct {
	User
	jwt.RegisteredClaims
}
