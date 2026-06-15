UPDATE users SET nick_name = COALESCE(given_name, '') WHERE nick_name IS NULL;
ALTER TABLE users ALTER COLUMN nick_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN nick_name SET DEFAULT '';
