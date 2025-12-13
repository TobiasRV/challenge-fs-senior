-- name: CreateProject :one
INSERT INTO projects (created_at, updated_at, name,team_id, manager_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET name = $1, status = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: GetProjectById :one
SELECT * FROM projects
WHERE id = $1
LIMIT 1;

-- name: GetProjectByManager :one
SELECT * FROM projects
WHERE manager_id = $1
LIMIT 1;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;