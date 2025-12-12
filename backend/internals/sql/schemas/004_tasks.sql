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
    name TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status TaskStatus NOT NULL DEFAULT 'ToDo',
    title TEXT NOT NULL,
    description TEXT
);

-- +goose Down
DROP TABLE tasks;