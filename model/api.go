// The package model contains structures that model API and user data.
package model

import "github.com/golang-jwt/jwt/v5"

// UserTokenResponse is returned when the user has made a successful request for a new login token.
type LoginTokenResponse struct {
	Token string `json:"token"`
}

// ResetTokenResponse is returned when the user has made a successful request for a new reset token.
type ResetTokenResponse struct {
	Token string `json:"reset_token"`
}

const (
	// This is a login token.
	TokenUsageLogin = "login"
	// This is a reset token.
	TokenUsageReset = "reset"
)

// CustomClaims represent claims that are specific to this microservice.
type CustomClaims struct {
	Usage string `json:"usage"`
}

// UserClaims are the registered claims for a user's login or reset token.
type UserClaims struct {
	CustomClaims
	User
	jwt.RegisteredClaims
}
