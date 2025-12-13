-- name: CreateTasks :one
INSERT INTO tasks (created_at, updated_at, project_id, title, description, user_id)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5
WHERE id = $6
RETURNING *;

-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1
LIMIT 1;


-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;