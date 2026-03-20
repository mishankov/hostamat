-- +migrate ID: add_provider_columns
-- +migrate Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50) NOT NULL DEFAULT 'legacy';
UPDATE users SET provider = 'legacy' WHERE provider IS NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider_user_id VARCHAR(255);

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS provider_user_id;
ALTER TABLE users DROP COLUMN IF EXISTS provider;
