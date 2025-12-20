-- name: CreateUserEvent :one
INSERT INTO user_events(request_id, event_seq, user_id, method, method_details, created_at)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        Now()
)
RETURNING request_id;

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
WHERE method = $1
ORDER BY  created_at DESC;

-- name: GetEventsByEndpoint :many
SELECT * FROM user_events
WHERE method_details::json->>'endpoint' = $1
ORDER BY created_at DESC;

-- name: CountEventsByUser :one
SELECT COUNT(*) FROM user_events 
WHERE user_id = $1;

-- name: CountEventsByAction :one
SELECT COUNT(*) FROM user_events
WHERE method = $1;

-- name: CountEventsByIP :one 
SELECT COUNT(*) FROM user_events
WHERE method_details::json->>'ip' = $1;

-- name: GetEventsByIP :many 
SELECT * FROM user_events
WHERE method_details::json->>'ip' = $1
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

-- name: GetState :one
SELECT state FROM user_events
WHERE request_id = $1;

-- name: SetStatePending :exec 
UPDATE user_events
SET state = 'Pending'
WHERE request_id = $1;

-- name: SetStateSuccess :exec
UPDATE user_events 
SET state = 'Success'
WHERE request_id = $1;

-- name: SetStateFailure :exec
UPDATE user_events
SET state = 'Failure'
WHERE request_id = $1;
