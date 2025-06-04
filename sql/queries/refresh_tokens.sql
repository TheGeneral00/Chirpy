-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
Values ( 
        $1, Now(), Now(), $2, $3, NULL
) RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = Now(), updated_at = Now()
WHERE token = $1;

-- name: RetrieveRefreshToken :one
SELECT token, expires_at, revoked_at, user_id FROM refresh_tokens
WHERE token = $1;
