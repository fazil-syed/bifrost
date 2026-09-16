-- +goose up

CREATE TABLE roles(
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL
        REFERENCES applications(id)
        ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    UNIQUE (application_id,name)
);

-- +goose down
DROP TABLE roles;
