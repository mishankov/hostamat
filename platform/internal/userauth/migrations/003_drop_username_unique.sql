-- +migrate ID: drop_username_unique
-- +migrate Up
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;

-- +migrate Down
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
