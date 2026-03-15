-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS subdomains (
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	"user" VARCHAR(255) NOT NULL,
	full_domain VARCHAR(255) NOT NULL
)

-- +migrate Down
DROP TABLE subdomains
