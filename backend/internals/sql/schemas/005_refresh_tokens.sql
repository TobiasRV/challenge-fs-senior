-- +goose Up
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    userId UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE refresh_tokens;