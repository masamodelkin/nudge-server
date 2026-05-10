-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT UNIQUE,
    auth_provider TEXT NOT NULL DEFAULT 'local',
    provider_id TEXT,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- +goose Down
DROP TABLE users;