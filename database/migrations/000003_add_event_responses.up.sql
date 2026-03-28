CREATE TABLE IF NOT EXISTS event_responses
(
    id         UUID         DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    event_id   UUID         REFERENCES events (id)    NOT NULL,
    user_id    UUID         REFERENCES users (id)     NOT NULL,
    response   BOOLEAN                                NOT NULL,
    created_at TIMESTAMP    DEFAULT current_timestamp NOT NULL,
    UNIQUE (event_id, user_id)
);