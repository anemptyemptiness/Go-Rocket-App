-- +goose Up
ALTER TABLE parts
    ADD COLUMN IF NOT EXISTS properties JSONB NOT NULL DEFAULT '{}';
ALTER TABLE parts
    ADD COLUMN IF NOT EXISTS reserved INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE parts
    DROP COLUMN IF EXISTS properties;
ALTER TABLE parts
    DROP COLUMN IF EXISTS reserved;