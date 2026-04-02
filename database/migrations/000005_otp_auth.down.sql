DROP TABLE IF EXISTS otp_tokens;

-- Restore issuer column (nullable — original data cannot be recovered)
ALTER TABLE users ADD COLUMN issuer VARCHAR(100);

-- Restore NOT NULL constraints
UPDATE users SET family_name = '' WHERE family_name IS NULL;
UPDATE users SET given_name = '' WHERE given_name IS NULL;
UPDATE users SET nick_name = '' WHERE nick_name IS NULL;
ALTER TABLE users ALTER COLUMN family_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN given_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN nick_name SET NOT NULL;
