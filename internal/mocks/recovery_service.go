package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AccountRecoveryService struct {
	mock.Mock
}

func (m *AccountRecoveryService) Start(ctx context.Context, phone string) error {
	args := m.Called(ctx, phone)
	return args.Error(0)
}

func (m *AccountRecoveryService) VerifyOTP(ctx context.Context, phone, code string) (string, error) {
	args := m.Called(ctx, phone, code)
	return args.String(0), args.Error(1)
}

func (m *AccountRecoveryService) ConfirmPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error) {
	args := m.Called(ctx, accountID, pin)
	return args.String(0), args.Error(1)
}

func (m *AccountRecoveryService) ResetPassword(ctx context.Context, accountID uuid.UUID, newPassword string) error {
	args := m.Called(ctx, accountID, newPassword)
	return args.Error(0)
}
