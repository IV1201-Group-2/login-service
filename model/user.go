package model

import "github.com/golang-jwt/jwt/v5"

type Role int

const (
	RoleRecruiter Role = iota + 1
	RoleApplicant
)

type User struct {
	ID   int  `json:"id"`
	Role Role `json:"role"`

	Username string `json:"username"`
	Email    string `json:"email"`

	// Omit from JSON response
	Password string `json:"-"`
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}
