-- +migrate ID: add_provider_user_id_unique_index
-- +migrate Up
CREATE UNIQUE INDEX IF NOT EXISTS users_provider_provider_user_id_idx ON users (provider, provider_user_id) WHERE provider_user_id IS NOT NULL;

-- +migrate Down
DROP INDEX IF EXISTS users_provider_provider_user_id_idx;
