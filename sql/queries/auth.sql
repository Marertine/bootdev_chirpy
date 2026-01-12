-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE token = $1;