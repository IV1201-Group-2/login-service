package model

const (
	APIErrUnknown           = "UNKNOWN"
	APIErrMissingParameters = "MISSING_PARAMETERS"

	APIErrWrongIdentity = "WRONG_IDENTITY"
	APIErrWrongPassword = "WRONG_PASSWORD"

	APIErrAlreadyLoggedIn = "ALREADY_LOGGED_IN"
)

type LoginSuccessResponse struct {
	Token string `json:"token"`
}

// Generic error response
type ErrorDetails struct {
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error   string        `json:"error"`
	Details *ErrorDetails `json:"details,omitempty"`
}
