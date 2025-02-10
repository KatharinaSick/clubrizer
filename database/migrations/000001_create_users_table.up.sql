CREATE TABLE IF NOT EXISTS users
(
    id          UUID DEFAULT gen_random_uuid() UNIQUE PRIMARY KEY,
    email       VARCHAR(300) UNIQUE NOT NULL,
    family_name VARCHAR(100)        NOT NULL,
    given_name  VARCHAR(100)        NOT NULL,
    nick_name   VARCHAR(100)        NOT NULL,
    picture     VARCHAR(500),
    issuer      VARCHAR(100)        NOT NULL
);