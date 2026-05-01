-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS api_keys (
	id VARCHAR(255) PRIMARY KEY,
	user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(80) NOT NULL,
	key_prefix VARCHAR(255) NOT NULL,
	secret_digest TEXT NOT NULL,
	scopes TEXT[] NOT NULL,
	expires_at TIMESTAMPTZ NULL,
	last_used_at TIMESTAMPTZ NULL,
	revoked_at TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_prefix ON api_keys(key_prefix);
CREATE INDEX IF NOT EXISTS idx_api_keys_user ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_active_user ON api_keys(user_id) WHERE revoked_at IS NULL;

-- +migrate Down
DROP TABLE IF EXISTS api_keys;
