-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, userId, token, expires_at, created_at, revoked)
VALUES($1, $2, $3, $4, $5, $6);

-- name: UpdaterefreshToken :exec
UPDATE refresh_tokens
SET token = $1, expires_at = $2, created_at = $3, revoked = $4
WHERE id = $5;

-- name: GetRefreshTokenByToken :one
SELECT refresh_tokens.*,
users.id AS userDataId,
users.created_at AS userDataCreatedAt,
users.updated_at AS userDataUpdatedAt,
users.username AS userDataUsername,
users.password AS userDataPassword,
users.email AS userDataEmail,
users.role AS userDataRole
FROM refresh_tokens
JOIN users ON refresh_tokens.userId = users.id
WHERE token = $1
LIMIT 1;

-- name: DeleteRefreshTokensByUserId :exec
DELETE FROM refresh_tokens
WHERE userId = $1;
