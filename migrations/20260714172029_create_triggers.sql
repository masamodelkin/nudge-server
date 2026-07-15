-- +goose Up
CREATE TABLE triggers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config TEXT NOT NULL,
    is_exclusive BOOLEAN NOT NULL DEFAULT false,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE task_triggers (
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    trigger_id TEXT NOT NULL REFERENCES triggers(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);

-- +goose Down
DROP TABLE triggers;
DROP TABLE task_triggers;