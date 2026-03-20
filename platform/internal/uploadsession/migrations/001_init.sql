-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS uploadsessions (
	id VARCHAR(255) PRIMARY KEY,
	"user" VARCHAR(255) NOT NULL,
	site VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	status VARCHAR(255) NOT NULL
)

-- +migrate Down
DROP TABLE uploadsessions
