// The package model contains structures that model API and user data.
package model

type UserTokenResponse struct {
	Token string `json:"token"`
}

type ResetTokenResponse struct {
	ResetToken string `json:"reset_token"`
}
