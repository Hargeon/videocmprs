-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS requests (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    user_id BIGINT REFERENCES users,
    video_name VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(255) NOT NULL DEFAULT 'original_in_review',
    details VARCHAR(255) NOT NULL DEFAULT '',

    bitrate BIGINT NOT NULL,
    resolution_x INT NOT NULL,
    resolution_y INT NOT NULL,
    ratio_x INT NOT NULL,
    ratio_y INT NOT NULL,

    original_file_id BIGINT,
    converted_file_id BIGINT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS requests;
-- +goose StatementEnd
