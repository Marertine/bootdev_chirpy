-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked_at)
VALUES ($1, $2, $3, NULL)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();


-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE token = $1;