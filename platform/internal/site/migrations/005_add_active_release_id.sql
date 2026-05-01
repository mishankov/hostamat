-- +migrate ID: add_active_release_id
-- +migrate Up
ALTER TABLE sites ADD COLUMN active_release_id VARCHAR(255)

-- +migrate Down
ALTER TABLE sites DROP COLUMN active_release_id
