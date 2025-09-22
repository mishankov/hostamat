package migrations

import "hostamat/pkg/platforma/database"

var Migrations = []database.Migration{
	{
		ID:   "add_user_column",
		Up:   `ALTER TABLE subdomains ADD COLUMN user VARCHAR(255)`,
		Down: `ALTER TABLE subdomains DROP COLUMN user`,
	},
}

var Schema = database.Schema{
	Queries: []string{
		`CREATE TABLE IF NOT EXISTS subdomains (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			user VARCHAR(255) NOT NULL,
			full_domain VARCHAR(255) NOT NULL
		)`,
	},
}
