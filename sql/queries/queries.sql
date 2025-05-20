-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (NOW(), NOW(), $1, $2) 
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (created_at, updated_at, body, user_id)
VALUES (NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;