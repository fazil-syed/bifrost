-- +goose up

CREATE TABLE external_identities (
    user_id UUID NOT NULL REFERENCES users(id),
    issuer TEXT NOT NULL,
    subject TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    
    PRIMARY KEY (issuer,subject)
);

-- +goose down
DROP TABLE external_identities;