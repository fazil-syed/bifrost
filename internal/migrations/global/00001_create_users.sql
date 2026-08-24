-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL
        CHECK (status in ('ACTIVE','DISABLED')),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE users;