-- +goose Up
CREATE TABLE user_events (
        id SERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        method TEXT NOT NULL,
        method_details TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        state TEXT DEFAULT 'Failed' CHECK (state IN('Failed', 'Success', 'Pending'))
);

-- +goose Down
DROP TABLE user_events;
