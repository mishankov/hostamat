-- +migrate ID: add_user_column
-- +migrate Up
ALTER TABLE sites ADD COLUMN "user" VARCHAR(255) NOT NULL

-- +migrate Down
ALTER TABLE sites DROP COLUMN "user"
