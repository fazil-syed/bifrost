-- +goose up

CREATE TABLE team_memberships(
    team_id UUID NOT NULL
        REFERENCES teams(id)
        ON DELETE CASCADE,

    user_id UUID NOT NULL,

    membership_type TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (team_id,user_id)
);

-- +goose down

DROP TABLE team_memberships;