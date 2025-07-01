-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserFromEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.*, refresh_tokens.expires_at as token_expires_at, refresh_tokens.revoked_at as token_revoked_at from users
join refresh_tokens on refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: UpdateUserEmailPassword :one
Update users set email = $1, hashed_password = $2 WHERE id = $3
RETURNING *;
