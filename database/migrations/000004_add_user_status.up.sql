CREATE TYPE user_status AS ENUM ('pending', 'approved', 'rejected');
ALTER TABLE users ADD COLUMN status user_status NOT NULL DEFAULT 'pending';
-- Existing users were active before approval was introduced, so grant them access automatically.
UPDATE users SET status = 'approved';
