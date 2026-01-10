-- name: CreateUser :one
INSERT INTO users (hashed_password, email)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one  
SELECT * FROM users
WHERE email = $1;