-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id bigserial NOT NULL unique PRIMARY KEY,
    email varchar(255) not null unique,
    password_hash varchar(255) not null,
    created_at  timestamp not null default current_timestamp()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
