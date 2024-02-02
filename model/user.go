package model

import "github.com/golang-jwt/jwt/v5"

// TODO: make sure this syncs up with the database

type Role int

const (
	RoleRecruiter Role = iota + 1
	RoleApplicant
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`

	Name           string `json:"name"`
	Surname        string `json:"surname"`
	PersonalNumber string `json:"pnr"`

	Role Role `json:"role"`
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}
