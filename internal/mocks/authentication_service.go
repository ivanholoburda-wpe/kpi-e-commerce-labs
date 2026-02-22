package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"lab2/internal/models"
)

type AuthenticationService struct {
	mock.Mock
}

func (m *AuthenticationService) Login(ctx context.Context, phone, password string) (string, error) {
	args := m.Called(ctx, phone, password)
	return args.String(0), args.Error(1)
}

func (m *AuthenticationService) Verify2FA(ctx context.Context, accountID uuid.UUID, code string) (string, string, error) {
	args := m.Called(ctx, accountID, code)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *AuthenticationService) SetupPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, string, error) {
	args := m.Called(ctx, accountID, pin)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *AuthenticationService) VerifyPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error) {
	args := m.Called(ctx, accountID, pin)
	return args.String(0), args.Error(1)
}

func (m *AuthenticationService) GetMe(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(models.Account), args.Error(1)
}
