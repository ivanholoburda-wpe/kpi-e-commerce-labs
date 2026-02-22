package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"lab2/internal/auth"
	"lab2/internal/models"
)

type AuthService struct {
	accounts   AccountRepository
	jwtSecret  []byte
	sessionTTL time.Duration // long-lived session token (days)
	accessTTL  time.Duration // short-lived access token (N hours, triggers PIN re-entry)
}

func NewAuthService(accounts AccountRepository, jwtSecret []byte, sessionTTL, accessTTL time.Duration) *AuthService {
	return &AuthService{
		accounts:   accounts,
		jwtSecret:  jwtSecret,
		sessionTTL: sessionTTL,
		accessTTL:  accessTTL,
	}
}

// Login verifies phone + password and returns a temporary JWT for 2FA verification.
func (s *AuthService) Login(ctx context.Context, phone, password string) (string, error) {
	account, err := s.accounts.GetByPhone(ctx, phone)
	if err != nil {
		return "", fmt.Errorf("get account: %w", err)
	}

	if account.PasswordHash == nil {
		return "", fmt.Errorf("account setup not complete")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*account.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := auth.GenerateToken(s.jwtSecret, account.ID, auth.Purpose2FA, 5*time.Minute)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// Verify2FA validates the TOTP code and returns:
//   - a session JWT + next_step "verify_pin"  — if PIN is already set
//   - a pin_setup JWT + next_step "setup_pin" — if PIN needs to be created
func (s *AuthService) Verify2FA(ctx context.Context, accountID uuid.UUID, code string) (token, nextStep string, err error) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return "", "", fmt.Errorf("get account: %w", err)
	}

	if account.TOTPSecret == nil {
		return "", "", fmt.Errorf("2fa not configured")
	}

	if !totp.Validate(code, *account.TOTPSecret) {
		return "", "", fmt.Errorf("invalid 2fa code")
	}

	if account.PinHash == nil {
		token, err = auth.GenerateToken(s.jwtSecret, accountID, auth.PurposePinSetup, 10*time.Minute)
		if err != nil {
			return "", "", fmt.Errorf("generate token: %w", err)
		}
		return token, "setup_pin", nil
	}

	token, err = auth.GenerateToken(s.jwtSecret, accountID, auth.PurposeSession, s.sessionTTL)
	if err != nil {
		return "", "", fmt.Errorf("generate token: %w", err)
	}
	return token, "verify_pin", nil
}

// SetupPIN hashes the PIN, saves it and returns a session + access token pair.
func (s *AuthService) SetupPIN(ctx context.Context, accountID uuid.UUID, pin string) (sessionToken, accessToken string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("hash pin: %w", err)
	}

	if err := s.accounts.UpdatePinHash(ctx, accountID, string(hash)); err != nil {
		return "", "", fmt.Errorf("save pin: %w", err)
	}

	sessionToken, err = auth.GenerateToken(s.jwtSecret, accountID, auth.PurposeSession, s.sessionTTL)
	if err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}

	accessToken, err = auth.GenerateToken(s.jwtSecret, accountID, auth.PurposeAccess, s.accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	return sessionToken, accessToken, nil
}

// VerifyPIN checks the PIN and returns a fresh access token (refreshes the N-hour window).
// Called with the long-lived session token.
func (s *AuthService) VerifyPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error) {
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

	accessToken, err := auth.GenerateToken(s.jwtSecret, accountID, auth.PurposeAccess, s.accessTTL)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, nil
}

// GetMe returns the account for the authenticated user.
func (s *AuthService) GetMe(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return models.Account{}, fmt.Errorf("get account: %w", err)
	}
	return account, nil
}
