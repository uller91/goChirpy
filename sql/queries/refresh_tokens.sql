-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: RevokeRefreshToken :one
Update refresh_tokens set revoked_at = NOW(), updated_at = NOW() WHERE token = $1
RETURNING *;
