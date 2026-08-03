-- Comments members can leave on an event, shown in the event detail view. A comment is
-- always written by an account holder (a user), never a kid. body is capped at 500 chars;
-- the application also validates this, but the column length is the hard backstop.
CREATE TABLE IF NOT EXISTS event_comments
(
    id         UUID         DEFAULT gen_random_uuid() PRIMARY KEY,
    event_id   UUID         REFERENCES events (id) ON DELETE CASCADE NOT NULL,
    user_id    UUID         REFERENCES users (id)                    NOT NULL,
    body       VARCHAR(500)                                          NOT NULL,
    created_at TIMESTAMP    DEFAULT current_timestamp                NOT NULL
);
