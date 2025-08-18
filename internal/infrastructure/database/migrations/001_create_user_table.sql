-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX email_index ON users (email);

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS email_index CASCADE;