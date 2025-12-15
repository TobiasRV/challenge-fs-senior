-- +goose Up
DROP TYPE IF EXISTS TaskStatus; CREATE TYPE TaskStatus AS ENUM (
  'ToDo',
  'InProgress',
  'Done'
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status TaskStatus NOT NULL DEFAULT 'ToDo',
    title TEXT NOT NULL,
    description TEXT
);

CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_tasks_title ON tasks(title);

-- +goose Down
DROP TABLE tasks;