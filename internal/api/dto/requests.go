package dto

// --- Registration ---

type StartRegistrationRequest struct {
	Phone          string `json:"phone" binding:"required"`
	TelegramChatID int64  `json:"telegram_chat_id" binding:"required"`
}

type VerifyCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type SetupRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type ConfirmTOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// --- Auth ---

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Verify2FARequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type SetupPINRequest struct {
	Pin string `json:"pin" binding:"required,min=4,max=6"`
}

type VerifyPINRequest struct {
	Pin string `json:"pin" binding:"required"`
}

// --- Recovery ---

type StartRecoveryRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyRecoveryOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type ConfirmRecoveryPINRequest struct {
	Pin string `json:"pin" binding:"required"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
