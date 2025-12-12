-- name: CreateTeam :one
INSERT INTO teams (created_at, updated_at, name, owner_id)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetTeamByOwner :one
SELECT * FROM teams
WHERE owner_id = $1
LIMIT 1;