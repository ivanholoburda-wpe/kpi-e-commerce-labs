-- name: CreateVerificationCode :one
INSERT INTO verification_codes (account_id, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetActiveCodeByAccountID :one
SELECT * FROM verification_codes
WHERE account_id = $1
  AND used_at IS NULL
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkCodeUsed :exec
UPDATE verification_codes SET used_at = now() WHERE id = $1;
