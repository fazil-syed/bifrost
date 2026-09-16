-- +goose up
CREATE TABLE tenants(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    database_name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    UNIQUE (slug),
    UNIQUE (database_name)

);

-- +goose down

DROP TABLE tenants;