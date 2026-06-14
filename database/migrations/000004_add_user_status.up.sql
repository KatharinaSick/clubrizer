CREATE TYPE IF NOT EXISTS user_status AS ENUM ('pending', 'approved', 'rejected');
-- DEFAULT 'approved' so existing users are grandfathered in; changed to 'pending' immediately after.
ALTER TABLE users ADD COLUMN IF NOT EXISTS status user_status NOT NULL DEFAULT 'approved';
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'pending';
