-- Kinds of push notifications a user switched off on the notification settings screen. Only
-- switched-off kinds are stored, so everyone without a row gets everything — new users and new
-- kinds start switched on without a backfill. Settings are per user, not per device.
CREATE TABLE IF NOT EXISTS notification_opt_outs
(
    user_id  UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    category TEXT                                         NOT NULL,
    PRIMARY KEY (user_id, category)
);
