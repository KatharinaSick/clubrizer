CREATE TYPE user_status AS ENUM ('pending', 'approved', 'rejected');
-- DEFAULT 'approved' so existing users are grandfathered in; changed to 'pending' immediately after.
ALTER TABLE users ADD COLUMN status user_status NOT NULL DEFAULT 'approved';
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'pending';
