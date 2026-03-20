-- +migrate ID: add_name_column
-- +migrate Up
ALTER TABLE sites ADD COLUMN "name" VARCHAR(255) NOT NULL DEFAULT 'Сайт'

-- +migrate Down
ALTER TABLE sites DROP COLUMN "name"
