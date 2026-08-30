-- +goose up
CREATE TABLE credentials (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- +goose down

DROP TABLE credentials;