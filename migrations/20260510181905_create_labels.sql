-- +goose Up
CREATE TABLE labels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    color TEXT,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE labels;