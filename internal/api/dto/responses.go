package dto

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type SetupResponse struct {
	TOTPUri string `json:"totp_uri"`
}

type Verify2FAResponse struct {
	Token    string `json:"token"`
	NextStep string `json:"next_step"` // "setup_pin" or "verify_pin"
}

type AuthTokensResponse struct {
	SessionToken string `json:"session_token"`
	AccessToken  string `json:"access_token"`
}

type AccountResponse struct {
	ID          string `json:"id"`
	Phone       string `json:"phone"`
	TOTPEnabled bool   `json:"totp_enabled"`
	PinSet      bool   `json:"pin_set"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
