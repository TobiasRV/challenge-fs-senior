-- +goose Up

DROP TYPE IF EXISTS UserRoles; CREATE TYPE UserRoles AS ENUM (
  'Admin',
  'Manager',
  'Member'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role UserRoles NOT NULL,
    team_id UUID
);

-- +goose Down
DROP TABLE users CASCADE;