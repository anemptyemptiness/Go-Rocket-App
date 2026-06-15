-- +goose Up
ALTER TABLE orders ADD COLUMN IF NOT EXISTS user_uuid UUID NOT NULL;

-- +goose Down
ALTER TABLE orders DROP COLUMN IF EXISTS user_uuid;