package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"lab2/internal/auth"
	"lab2/internal/models"
)

type RegisterService struct {
	accounts  AccountRepository
	codes     VerificationCodeRepository
	notifier  Notifier
	jwtSecret []byte
}

func NewRegisterService(
	accounts AccountRepository,
	codes VerificationCodeRepository,
	notifier Notifier,
	jwtSecret []byte,
) *RegisterService {
	return &RegisterService{
		accounts:  accounts,
		codes:     codes,
		notifier:  notifier,
		jwtSecret: jwtSecret,
	}
}

// Start creates an account and sends an OTP code via the notifier.
func (s *RegisterService) Start(ctx context.Context, phone string, telegramChatID int64) error {
	account, err := s.accounts.Create(ctx, phone, telegramChatID)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	if _, err := s.codes.Create(ctx, account.ID, code, expiresAt); err != nil {
		return fmt.Errorf("save verification code: %w", err)
	}

	message := fmt.Sprintf("Your verification code: %s", code)
	if err := s.notifier.SendMessage(ctx, telegramChatID, message); err != nil {
		return fmt.Errorf("send otp: %w", err)
	}

	return nil
}

// Verify checks the OTP and returns a scoped JWT for the setup step.
func (s *RegisterService) Verify(ctx context.Context, phone string, code string) (string, error) {
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

	if err := s.accounts.UpdateStatus(ctx, account.ID, models.StatusPendingSetup); err != nil {
		return "", fmt.Errorf("update status: %w", err)
	}

	token, err := auth.GenerateToken(s.jwtSecret, account.ID, auth.PurposeRegistration, 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// Setup hashes the password, saves it, generates a TOTP secret and returns the provisioning URI.
func (s *RegisterService) Setup(ctx context.Context, accountID uuid.UUID, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	if err := s.accounts.UpdatePasswordHash(ctx, accountID, string(hash)); err != nil {
		return "", fmt.Errorf("save password: %w", err)
	}

	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("get account: %w", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Lab2",
		AccountName: account.Phone,
	})
	if err != nil {
		return "", fmt.Errorf("generate totp: %w", err)
	}

	if err := s.accounts.UpdateTOTPSecret(ctx, accountID, key.Secret()); err != nil {
		return "", fmt.Errorf("save totp secret: %w", err)
	}

	return key.URL(), nil
}

// ConfirmTOTP verifies the TOTP code from the authenticator app and activates the account.
func (s *RegisterService) ConfirmTOTP(ctx context.Context, accountID uuid.UUID, code string) error {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	if account.TOTPSecret == nil {
		return fmt.Errorf("totp not configured")
	}

	valid := totp.Validate(code, *account.TOTPSecret)
	if !valid {
		return fmt.Errorf("invalid totp code")
	}

	if err := s.accounts.EnableTOTP(ctx, accountID); err != nil {
		return fmt.Errorf("enable totp: %w", err)
	}

	return nil
}

func generateOTP(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", length, n), nil
}
