-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS sites (
	id VARCHAR(255) PRIMARY KEY,
	subdomain VARCHAR(255) NOT NULL,
	domain VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	version INT NOT NULL
)

-- +migrate Down
DROP TABLE IF EXISTS sites
