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

CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_name ON projects(name);

-- +goose Down
DROP TABLE projects;