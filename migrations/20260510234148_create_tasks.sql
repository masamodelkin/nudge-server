-- +goose Up
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    is_draft BOOLEAN NOT NULL DEFAULT false,
    due_date INTEGER,
    priority INTEGER,
    duration INTEGER,
    time_spent INTEGER NOT NULL DEFAULT 0,
    status_id TEXT REFERENCES statuses(id),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- +goose Down
DROP TABLE tasks;