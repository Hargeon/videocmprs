-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS videos (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    bitrate BIGINT,
    resolution varchar(255),
    ratio varchar(255),
    service_id varchar(255),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos;
-- +goose StatementEnd
