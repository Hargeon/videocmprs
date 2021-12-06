-- +goose Up
CREATE TABLE IF NOT EXISTS videos (
    id BIGSERIAL NOT NULL UNIQUE PRIMARY KEY,
    user_id BIGINT REFERENCES users,
    name VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,

    bitrate BIGINT,
    resolution_x INT,
    resolution_y INT,
    ratio_x INT,
    ratio_y INT,

    service_id varchar(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS videos;
