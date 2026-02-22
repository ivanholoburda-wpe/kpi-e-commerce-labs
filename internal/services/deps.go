package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"lab2/internal/models"
)

// AccountRepository is the shared data-access contract for accounts,
// used by registration, authentication and recovery services.
type AccountRepository interface {
	Create(ctx context.Context, phone string, telegramChatID int64) (models.Account, error)
	GetByPhone(ctx context.Context, phone string) (models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Account, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error
	UpdateTOTPSecret(ctx context.Context, id uuid.UUID, secret string) error
	EnableTOTP(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.AccountStatus) error
	UpdatePinHash(ctx context.Context, id uuid.UUID, hash string) error
}

// VerificationCodeRepository is the shared data-access contract for OTP codes.
type VerificationCodeRepository interface {
	Create(ctx context.Context, accountID uuid.UUID, code string, expiresAt time.Time) (models.VerificationCode, error)
	GetActiveByAccountID(ctx context.Context, accountID uuid.UUID) (models.VerificationCode, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
}

// Notifier sends messages to users (e.g. via Telegram).
type Notifier interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
