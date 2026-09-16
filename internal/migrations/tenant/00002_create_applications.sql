-- +goose up

CREATE TABLE applications(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    owner_user_id UUID,
    owner_team_id UUID
        REFERENCES teams(id)
        ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    UNIQUE (slug)
);


-- +goose down

DROP TABLE applications;