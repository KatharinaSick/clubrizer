-- Name fields are now filled in during profile setup (after admin approval), not at account creation
ALTER TABLE users ALTER COLUMN family_name DROP NOT NULL;
ALTER TABLE users ALTER COLUMN given_name DROP NOT NULL;
ALTER TABLE users ALTER COLUMN nick_name DROP NOT NULL;

-- Drop Google/OAuth-specific column
ALTER TABLE users DROP COLUMN issuer;

CREATE TABLE IF NOT EXISTS otp_tokens
(
    id             UUID        DEFAULT gen_random_uuid() PRIMARY KEY,
    email          VARCHAR(300)              NOT NULL,
    code_hash      VARCHAR(64)               NOT NULL,
    expires_at     TIMESTAMPTZ               NOT NULL,
    attempt_count  INT         DEFAULT 0    NOT NULL,
    invalidated_at TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_otp_tokens_email ON otp_tokens (email);
