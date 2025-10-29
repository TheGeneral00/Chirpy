-- +goose Up
CREATE TABLE user_events (
        id SERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        methode TEXT NOT NULL,
        methode_details TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE user_events;
