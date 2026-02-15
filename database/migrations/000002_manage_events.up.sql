CREATE TABLE IF NOT EXISTS event_categories
(
    id         UUID      DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    name       VARCHAR(50) UNIQUE                  NOT NULL,
    color      VARCHAR(7)                          NOT NULL,
    picture    VARCHAR(500)                        NOT NULL,
    sort_order INT       DEFAULT 0                 NOT NULL,
    created_by UUID REFERENCES users (id)          NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS events
(
    id          UUID      DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    title       VARCHAR(300)                          NOT NULL,
    category    UUID REFERENCES event_categories (id) NOT NULL,
    description VARCHAR(500),
    location    VARCHAR(500),
    start_time  TIMESTAMP                             NOT NULL,
    created_by  UUID REFERENCES users (id)            NOT NULL,
    created_at  TIMESTAMP DEFAULT current_timestamp   NOT NULL
);