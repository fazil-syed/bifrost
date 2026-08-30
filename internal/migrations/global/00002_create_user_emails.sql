-- +goose up 

CREATE TABLE user_emails (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    email TEXT NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);


-- +goose down

DROP TABLE user_emails;