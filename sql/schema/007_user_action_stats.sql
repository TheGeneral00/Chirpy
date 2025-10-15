-- +goose Up
ALTER TABLE users
ADD COLUMN action_count INT DEFAULT 0,
ADD COLUMN last_action TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN action_count,
DROP COLUMN last_action;
