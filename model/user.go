package model

import "github.com/golang-jwt/jwt/v5"

type Role int

const (
	RoleRecruiter Role = iota + 1
	RoleApplicant
)

type User struct {
	ID int `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`

	// Omit from JSON response
	Password string `json:"-"`

	Name           string `json:"name"`
	Surname        string `json:"surname"`
	PersonalNumber string `json:"pnr"`

	Role Role `json:"role"`
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}
