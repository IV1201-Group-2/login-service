// The package model contains structures that model API and user data.
package model

const (
	// Unknown error.
	APIErrUnknown = "UNKNOWN"

	// User did not provide identity, password or desired role.
	APIErrMissingParameters = "MISSING_PARAMETERS"
	// User does not have a password in the database.
	APIErrMissingPassword = "MISSING_PASSWORD"

	// No account was found with that specific username or email address.
	APIErrWrongIdentity = "WRONG_IDENTITY"
	// Account was found but the wrong password was provided.
	APIErrWrongPassword = "WRONG_PASSWORD"

	// User is already logged in (JWT token was provided).
	APIErrAlreadyLoggedIn = "ALREADY_LOGGED_IN"
)

// Specific success response for this API.
type SuccessResponse struct {
	Token string `json:"token"`
}

// Provides additional details about an error.
type ErrorDetails struct {
	Message string `json:"message"`
}

// Shared error response for all APIs.
type ErrorResponse struct {
	Error   string        `json:"error"`
	Details *ErrorDetails `json:"details,omitempty"`
}
