-- Browsers/devices that agreed to receive push notifications. Each row is one Web Push
-- subscription: the push service URL (endpoint) plus the keys used to encrypt messages for
-- that browser. A user can have several (phone, laptop). The endpoint is unique per
-- browser, so re-subscribing the same browser updates its row instead of adding a new one.
CREATE TABLE IF NOT EXISTS push_subscriptions
(
    id           UUID      DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id      UUID      REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    endpoint     TEXT      UNIQUE                                  NOT NULL,
    p256dh       TEXT                                              NOT NULL,
    auth         TEXT                                              NOT NULL,
    created_at   TIMESTAMP DEFAULT current_timestamp               NOT NULL,
    last_used_at TIMESTAMP DEFAULT current_timestamp               NOT NULL
);

CREATE INDEX IF NOT EXISTS push_subscriptions_user_id_idx ON push_subscriptions (user_id);
