-- +goose Up
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Add the fk from users after creating the table to circunvent the circular dependency
ALTER TABLE users
    ADD CONSTRAINT fk_users_teams
    FOREIGN KEY (team_id) REFERENCES teams(id);
-- +goose Down
DROP TABLE teams CASCADE;