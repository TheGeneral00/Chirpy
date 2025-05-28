-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
Values (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1, 
        $2
)
RETURNING *;

-- name: DeleteChrip :one
DELETE FROM chirps WHERE id = $1
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: ResetChirps :exec
TRUNCATE TABLE chirps;
