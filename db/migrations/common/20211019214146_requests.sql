-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS requests (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    user_id BIGINT REFERENCES users,
    status VARCHAR(255),
    details VARCHAR(255),

    bitrate BIGINT,
    resolution varchar(255),
    ratio varchar(255)

    original_file_id BIGINT,
    converted_file_id BIGINT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
-- +goose StatementEnd
