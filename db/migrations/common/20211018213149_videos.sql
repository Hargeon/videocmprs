-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS videos (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size BIGINT,
    bitrate BIGINT,
    resolution varchar(255),
    service_id varchar(255),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos
-- +goose StatementEnd
