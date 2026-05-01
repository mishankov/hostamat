-- +migrate ID: add_email
-- +migrate Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS email;
