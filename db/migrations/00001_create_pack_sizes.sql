-- +goose Up
CREATE TABLE pack_sizes (
    id         SERIAL PRIMARY KEY,
    size       SMALLINT NOT NULL UNIQUE CHECK (size > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE pack_sizes;
