-- +goose Up
ALTER TABLE statuses ADD COLUMN next_status_id TEXT REFERENCES statuses(id) ON DELETE SET NULL;
ALTER TABLE statuses ADD COLUMN is_done BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE statuses DROP COLUMN next_status_id;
ALTER TABLE statuses DROP COLUMN is_done;