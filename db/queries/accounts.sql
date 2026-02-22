-- name: CreateAccount :one
INSERT INTO accounts (phone, telegram_chat_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccountByPhone :one
SELECT * FROM accounts WHERE phone = $1;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: UpdatePasswordHash :exec
UPDATE accounts SET password_hash = $2, updated_at = now() WHERE id = $1;

-- name: UpdateTOTPSecret :exec
UPDATE accounts SET totp_secret = $2, updated_at = now() WHERE id = $1;

-- name: EnableTOTP :exec
UPDATE accounts SET totp_enabled = TRUE, status = 'active', updated_at = now() WHERE id = $1;

-- name: UpdateStatus :exec
UPDATE accounts SET status = $2, updated_at = now() WHERE id = $1;

-- name: UpdatePinHash :exec
UPDATE accounts SET pin_hash = $2, updated_at = now() WHERE id = $1;
