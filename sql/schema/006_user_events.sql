-- +goose Up
CREATE TABLE user_events (
        request_id UUID NOT NULL,
        event_seq INT NOT NULL,
        user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
        method TEXT NOT NULL,
        method_details TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        state TEXT DEFAULT 'Failed' CHECK (state IN('Failed', 'Success', 'Pending')),
        PRIMARY KEY (request_id, event_seq)
);

-- +goose Down
DROP TABLE user_events;
