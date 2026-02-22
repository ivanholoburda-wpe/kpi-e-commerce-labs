package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RegistrationService struct {
	mock.Mock
}

func (m *RegistrationService) Start(ctx context.Context, phone string, telegramChatID int64) error {
	args := m.Called(ctx, phone, telegramChatID)
	return args.Error(0)
}

func (m *RegistrationService) Verify(ctx context.Context, phone string, code string) (string, error) {
	args := m.Called(ctx, phone, code)
	return args.String(0), args.Error(1)
}

func (m *RegistrationService) Setup(ctx context.Context, accountID uuid.UUID, password string) (string, error) {
	args := m.Called(ctx, accountID, password)
	return args.String(0), args.Error(1)
}

func (m *RegistrationService) ConfirmTOTP(ctx context.Context, accountID uuid.UUID, code string) error {
	args := m.Called(ctx, accountID, code)
	return args.Error(0)
}
