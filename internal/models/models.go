package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	StatusPendingVerification AccountStatus = "pending_verification"
	StatusPendingSetup        AccountStatus = "pending_setup"
	StatusActive              AccountStatus = "active"
	StatusBlocked             AccountStatus = "blocked"
)

type Account struct {
	ID             uuid.UUID
	Phone          string
	TelegramChatID int64
	PasswordHash   *string
	TOTPSecret     *string
	TOTPEnabled    bool
	PinHash        *string
	Status         AccountStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type VerificationCode struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Code        string
	Attempts    int32
	MaxAttempts int32
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   time.Time
}
