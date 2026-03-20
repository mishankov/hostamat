-- +migrate ID: drop_version_column
-- +migrate Up
ALTER TABLE sites DROP COLUMN "version"

-- +migrate Down
ALTER TABLE sites ADD COLUMN "version" INT
