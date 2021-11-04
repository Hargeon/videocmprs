-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial NOT NULL UNIQUE PRIMARY KEY,
    email VARCHAR(255) NOT NULL unique,
    password_hash VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
