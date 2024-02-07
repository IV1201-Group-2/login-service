package model

const (
	APIErrUnknown = "UNKNOWN"

	APIErrMissingParameters = "MISSING_PARAMETERS"
	APIErrMissingPassword   = "MISSING_PASSWORD"

	APIErrWrongIdentity = "WRONG_IDENTITY"
	APIErrWrongPassword = "WRONG_PASSWORD"

	APIErrAlreadyLoggedIn = "ALREADY_LOGGED_IN"
)

// Specific success response for this API
type SuccessResponse struct {
	Token string `json:"token"`
}

// Shared error response for all APIs
type ErrorDetails struct {
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error   string        `json:"error"`
	Details *ErrorDetails `json:"details,omitempty"`
}
