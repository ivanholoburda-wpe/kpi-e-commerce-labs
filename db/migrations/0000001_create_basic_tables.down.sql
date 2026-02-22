DROP INDEX IF EXISTS idx_verification_codes_expires_at;
DROP INDEX IF EXISTS idx_verification_codes_account_id;
DROP TABLE IF EXISTS verification_codes;

DROP INDEX IF EXISTS idx_accounts_status;
DROP INDEX IF EXISTS idx_accounts_phone;
DROP TABLE IF EXISTS accounts;

DROP TYPE IF EXISTS account_status;
