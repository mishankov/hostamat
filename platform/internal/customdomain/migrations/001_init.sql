-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS custom_domains (
	id VARCHAR(255) PRIMARY KEY,
	site VARCHAR(255) NOT NULL,
	"user" VARCHAR(255) NOT NULL,
	hostname VARCHAR(255) NOT NULL,
	domain_type VARCHAR(32) NOT NULL,
	status VARCHAR(32) NOT NULL,
	verification_token VARCHAR(255) NOT NULL,
	last_verification_error TEXT NOT NULL DEFAULT '',
	last_checked_at TIMESTAMPTZ NULL,
	activated_at TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_custom_domains_user ON custom_domains("user");
CREATE INDEX IF NOT EXISTS idx_custom_domains_site ON custom_domains(site);
CREATE INDEX IF NOT EXISTS idx_custom_domains_status ON custom_domains(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_custom_domains_active_hostname ON custom_domains(hostname) WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS custom_domains;
