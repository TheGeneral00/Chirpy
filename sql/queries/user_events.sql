-- name: CreateUserEvent :one
INSERT INTO user_events(user_id, action, action_details, created_at)
VALUES (
        $1,
        $2,
        $3,
        Now()
)
RETURNING *;

-- name: ResetEvents :exec
TRUNCATE TABLE user_events RESTART IDENTITY;

-- name: GetEvents :many
SELECT * FROM user_events
ORDER BY created_at DESC;

-- name: GetEventCount :one
SELECT COUNT(*) FROM user_events;

-- name: GetEventsInTimeWindow :many
SELECT * FROM user_events
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC;

-- name: GetEventsByUser :many
SELECT *  FROM user_events
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetEventsByAction :many
SELECT * FROM user_events
WHERE action = $1
ORDER BY  created_at DESC;

-- name: GetEventsByEndpoint :many
SELECT * FROM user_events
WHERE action_details::json->>'endpoint' = $1
ORDER BY created_at DESC;

-- name: CountEventsByUser :one
SELECT COUNT(*) FROM user_events 
WHERE user_id = $1;

-- name: CountEventsByAction :one
SELECT COUNT(*) FROM user_events
WHERE action = $1;

-- name: CountEventsByIP :one 
SELECT COUNT(*) FROM user_events
WHERE action_details::json->>'ip' = $1;

-- name: GetEventsByIP :many 
SELECT * FROM user_events
WHERE action_details::json->>'ip' = $1
ORDER BY created_at DESC;

-- name: GetLatestEvents :many
SELECT * FROM user_events
ORDER BY created_at DESC
LIMIT $1;

-- name: GetLatestEventsByUser :many
SELECT * FROM user_events
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;
