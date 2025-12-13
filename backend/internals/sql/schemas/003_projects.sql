-- +goose Up

DROP TYPE IF EXISTS ProjectStatus; CREATE TYPE ProjectStatus AS ENUM (
  'OnHold',
  'InProgress',
  'Completed'
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    manager_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status ProjectStatus NOT NULL DEFAULT 'OnHold'
    
);

-- +goose Down
DROP TABLE projects;