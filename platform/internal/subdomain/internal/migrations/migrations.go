package migrations

import "github.com/mishankov/platforma/database"

var Migrations = []database.Migration{
	{
		ID: "init",
		Up: `CREATE TABLE IF NOT EXISTS subdomains (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			"user" VARCHAR(255) NOT NULL,
			full_domain VARCHAR(255) NOT NULL
		)`,
		Down: `DROP TABLE subdomains`,
	},
}
