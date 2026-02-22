CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE account_status AS ENUM (
    'pending_verification',
    'pending_setup',
    'active',
    'blocked'
);

CREATE TABLE accounts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone            VARCHAR(20)  NOT NULL UNIQUE,
    telegram_chat_id BIGINT       NOT NULL,
    password_hash VARCHAR(255),
    totp_secret   VARCHAR(128),
    totp_enabled  BOOLEAN      NOT NULL DEFAULT FALSE,
    status        account_status NOT NULL DEFAULT 'pending_verification',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_accounts_phone ON accounts (phone);
CREATE INDEX idx_accounts_status ON accounts (status);

CREATE TABLE verification_codes (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   UUID         NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    code         VARCHAR(10)  NOT NULL,
    attempts     INT          NOT NULL DEFAULT 0,
    max_attempts INT          NOT NULL DEFAULT 5,
    expires_at   TIMESTAMPTZ  NOT NULL,
    used_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_verification_codes_account_id ON verification_codes (account_id);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes (expires_at)
    WHERE used_at IS NULL;
