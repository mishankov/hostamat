-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(255) PRIMARY KEY,
	username VARCHAR(255),
	password TEXT,
	salt TEXT,
	created TIMESTAMP,
	updated TIMESTAMP,
	status VARCHAR(50),
	provider VARCHAR(50) NOT NULL DEFAULT 'legacy',
	provider_user_id VARCHAR(255)
)

-- +migrate Down
DROP TABLE users
