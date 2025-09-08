package migrations

import "hostamat/pkg/platforma/database"

var Migrations = []database.Migration{}

var Schema = database.Schema{
	Queries: []string{
		`CREATE TABLE IF NOT EXISTS subdomains (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			full_domain VARCHAR(255) NOT NULL
		)`,
	},
}
