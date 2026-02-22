package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lab2/internal/auth"
)

type RecoveryService struct {
	accounts  AccountRepository
	codes     VerificationCodeRepository
	notifier  Notifier
	jwtSecret []byte
}

func NewRecoveryService(
	accounts AccountRepository,
	codes VerificationCodeRepository,
	notifier Notifier,
	jwtSecret []byte,
) *RecoveryService {
	return &RecoveryService{
		accounts:  accounts,
		codes:     codes,
		notifier:  notifier,
		jwtSecret: jwtSecret,
	}
}

// Start sends an OTP to the user's Telegram for account recovery.
func (s *RecoveryService) Start(ctx context.Context, phone string) error {
	account, err := s.accounts.GetByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	if _, err := s.codes.Create(ctx, account.ID, code, expiresAt); err != nil {
		return fmt.Errorf("save verification code: %w", err)
	}

	message := fmt.Sprintf("Your recovery code: %s", code)
	if err := s.notifier.SendMessage(ctx, account.TelegramChatID, message); err != nil {
		return fmt.Errorf("send otp: %w", err)
	}

	return nil
}

// VerifyOTP checks the recovery OTP and returns a scoped recovery JWT.
func (s *RecoveryService) VerifyOTP(ctx context.Context, phone, code string) (string, error) {
	account, err := s.accounts.GetByPhone(ctx, phone)
	if err != nil {
		return "", fmt.Errorf("get account: %w", err)
	}

	vc, err := s.codes.GetActiveByAccountID(ctx, account.ID)
	if err != nil {
		return "", fmt.Errorf("get verification code: %w", err)
	}

	if vc.Code != code {
		return "", fmt.Errorf("invalid verification code")
	}

	if err := s.codes.MarkUsed(ctx, vc.ID); err != nil {
		return "", fmt.Errorf("mark code used: %w", err)
	}

	token, err := auth.GenerateToken(s.jwtSecret, account.ID, auth.PurposeRecoveryOTP, 10*time.Minute)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// ConfirmPIN verifies the user's existing PIN as an additional recovery factor.
func (s *RecoveryService) ConfirmPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("get account: %w", err)
	}

	if account.PinHash == nil {
		return "", fmt.Errorf("pin not set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*account.PinHash), []byte(pin)); err != nil {
		return "", fmt.Errorf("invalid pin")
	}

	token, err := auth.GenerateToken(s.jwtSecret, accountID, auth.PurposeRecoveryConfirmed, 10*time.Minute)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// ResetPassword sets a new password for the account.
func (s *RecoveryService) ResetPassword(ctx context.Context, accountID uuid.UUID, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.accounts.UpdatePasswordHash(ctx, accountID, string(hash)); err != nil {
		return fmt.Errorf("save password: %w", err)
	}

	return nil
}
