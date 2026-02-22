package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"lab2/internal/db/sqlc"
	"lab2/internal/models"
)

type AccountRepo struct {
	q *sqlc.Queries
}

func NewAccountRepo(q *sqlc.Queries) *AccountRepo {
	return &AccountRepo{q: q}
}

func (r *AccountRepo) Create(ctx context.Context, phone string, telegramChatID int64) (models.Account, error) {
	row, err := r.q.CreateAccount(ctx, sqlc.CreateAccountParams{
		Phone:          phone,
		TelegramChatID: telegramChatID,
	})
	if err != nil {
		return models.Account{}, err
	}
	return toAccount(row), nil
}

func (r *AccountRepo) GetByPhone(ctx context.Context, phone string) (models.Account, error) {
	row, err := r.q.GetAccountByPhone(ctx, phone)
	if err != nil {
		return models.Account{}, err
	}
	return toAccount(row), nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	row, err := r.q.GetAccountByID(ctx, id)
	if err != nil {
		return models.Account{}, err
	}
	return toAccount(row), nil
}

func (r *AccountRepo) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	return r.q.UpdatePasswordHash(ctx, sqlc.UpdatePasswordHashParams{
		ID:           id,
		PasswordHash: pgtype.Text{String: hash, Valid: true},
	})
}

func (r *AccountRepo) UpdateTOTPSecret(ctx context.Context, id uuid.UUID, secret string) error {
	return r.q.UpdateTOTPSecret(ctx, sqlc.UpdateTOTPSecretParams{
		ID:         id,
		TotpSecret: pgtype.Text{String: secret, Valid: true},
	})
}

func (r *AccountRepo) EnableTOTP(ctx context.Context, id uuid.UUID) error {
	return r.q.EnableTOTP(ctx, id)
}

func (r *AccountRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.AccountStatus) error {
	return r.q.UpdateStatus(ctx, sqlc.UpdateStatusParams{
		ID:     id,
		Status: sqlc.AccountStatus(status),
	})
}

func (r *AccountRepo) UpdatePinHash(ctx context.Context, id uuid.UUID, hash string) error {
	return r.q.UpdatePinHash(ctx, sqlc.UpdatePinHashParams{
		ID:      id,
		PinHash: pgtype.Text{String: hash, Valid: true},
	})
}

func toAccount(row sqlc.Account) models.Account {
	var passwordHash *string
	if row.PasswordHash.Valid {
		passwordHash = &row.PasswordHash.String
	}

	var totpSecret *string
	if row.TotpSecret.Valid {
		totpSecret = &row.TotpSecret.String
	}

	var pinHash *string
	if row.PinHash.Valid {
		pinHash = &row.PinHash.String
	}

	return models.Account{
		ID:             row.ID,
		Phone:          row.Phone,
		TelegramChatID: row.TelegramChatID,
		PasswordHash:   passwordHash,
		TOTPSecret:     totpSecret,
		TOTPEnabled:    row.TotpEnabled,
		PinHash:        pinHash,
		Status:         models.AccountStatus(row.Status),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}
