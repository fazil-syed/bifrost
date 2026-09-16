-- +goose up

CREATE TABLE role_assignments(
    id UUID PRIMARY KEY,

    role_id UUID NOT NULL
        REFERENCES roles(id)
        ON DELETE CASCADE,
    user_id UUID,
    team_id UUID
        REFERENCES teams(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX role_assignments_user_unique
    ON role_assignments (role_id,user_id)
    WHERE user_id IS NOT NULL;


CREATE UNIQUE INDEX role_assignments_team_unique
    ON role_assignments (role_id,team_id)
    WHERE team_id IS NOT NULL;


-- +goose down

DROP TABLE role_assignments;