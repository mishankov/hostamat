-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS form_discoveries (
	id VARCHAR(255) PRIMARY KEY,
	site VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_form_discoveries_site ON form_discoveries(site);

CREATE TABLE IF NOT EXISTS form_definitions (
	id VARCHAR(255) PRIMARY KEY,
	discovery_id VARCHAR(255) NOT NULL,
	site VARCHAR(255) NOT NULL,
	file_name TEXT NOT NULL,
	page_path TEXT NOT NULL,
	form_key VARCHAR(128) NOT NULL,
	label TEXT NOT NULL,
	original_action TEXT NOT NULL,
	method VARCHAR(16) NOT NULL,
	fields JSONB NOT NULL,
	supported BOOLEAN NOT NULL,
	unsupported_reason TEXT NOT NULL,
	active BOOLEAN NOT NULL,
	public_token VARCHAR(255) NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_form_definitions_site ON form_definitions(site);
CREATE INDEX IF NOT EXISTS idx_form_definitions_discovery ON form_definitions(discovery_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_form_definitions_public_token ON form_definitions(public_token) WHERE public_token <> '';

CREATE TABLE IF NOT EXISTS form_submissions (
	id VARCHAR(255) PRIMARY KEY,
	site VARCHAR(255) NOT NULL,
	form VARCHAR(255) NOT NULL,
	page_path TEXT NOT NULL,
	fields JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_form_submissions_site_created ON form_submissions(site, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_form_submissions_form_created ON form_submissions(form, created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS form_submissions;
DROP TABLE IF EXISTS form_definitions;
DROP TABLE IF EXISTS form_discoveries;
