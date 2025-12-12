-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, username, password, email, role)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;
