-- +migrate ID: add_site_warnings
-- +migrate Up
CREATE TABLE IF NOT EXISTS site_warnings (
	site VARCHAR(255) NOT NULL,
	code VARCHAR(64) NOT NULL,
	severity VARCHAR(32) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	PRIMARY KEY (site, code)
);

CREATE INDEX IF NOT EXISTS idx_site_warnings_site ON site_warnings(site);

-- +migrate Down
DROP TABLE IF EXISTS site_warnings;
