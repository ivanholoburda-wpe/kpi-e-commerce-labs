package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"lab2/internal/db/sqlc"
	"lab2/internal/models"
)

type VerificationCodeRepo struct {
	q *sqlc.Queries
}

func NewVerificationCodeRepo(q *sqlc.Queries) *VerificationCodeRepo {
	return &VerificationCodeRepo{q: q}
}

func (r *VerificationCodeRepo) Create(ctx context.Context, accountID uuid.UUID, code string, expiresAt time.Time) (models.VerificationCode, error) {
	row, err := r.q.CreateVerificationCode(ctx, sqlc.CreateVerificationCodeParams{
		AccountID: accountID,
		Code:      code,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return models.VerificationCode{}, err
	}
	return toVerificationCode(row), nil
}

func (r *VerificationCodeRepo) GetActiveByAccountID(ctx context.Context, accountID uuid.UUID) (models.VerificationCode, error) {
	row, err := r.q.GetActiveCodeByAccountID(ctx, accountID)
	if err != nil {
		return models.VerificationCode{}, err
	}
	return toVerificationCode(row), nil
}

func (r *VerificationCodeRepo) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkCodeUsed(ctx, id)
}

func toVerificationCode(row sqlc.VerificationCode) models.VerificationCode {
	var usedAt *time.Time
	if row.UsedAt.Valid {
		usedAt = &row.UsedAt.Time
	}

	return models.VerificationCode{
		ID:          row.ID,
		AccountID:   row.AccountID,
		Code:        row.Code,
		Attempts:    row.Attempts,
		MaxAttempts: row.MaxAttempts,
		ExpiresAt:   row.ExpiresAt.Time,
		UsedAt:      usedAt,
		CreatedAt:   row.CreatedAt.Time,
	}
}
