-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS requests (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    user_id BIGINT REFERENCES users,
    status VARCHAR(255) NOT NULL DEFAULT 'original_in_review',
    details VARCHAR(255) NOT NULL DEFAULT '',

    bitrate BIGINT NOT NULL,
    resolution varchar(255) NOT NULL,
    ratio varchar(255) NOT NULL,

    original_file_id BIGINT,
    converted_file_id BIGINT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
-- +goose StatementEnd
