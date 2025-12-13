-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, username, password, email, role, team_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET username = $1, email = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;