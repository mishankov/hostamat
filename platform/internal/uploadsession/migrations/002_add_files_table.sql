-- +migrate ID: add_files_table
-- +migrate Up
CREATE TABLE IF NOT EXISTS files (
	id VARCHAR(255) PRIMARY KEY,
	upload_session VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	uploaded BOOLEAN NOT NULL,
	size INT NOT NULL
)

-- +migrate Down
DROP TABLE files
