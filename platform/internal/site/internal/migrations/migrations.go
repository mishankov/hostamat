package migrations

import "github.com/platforma-dev/platforma/database"

var Migrations = []database.Migration{
	{
		ID: "init",
		Up: `CREATE TABLE IF NOT EXISTS sites (
					id VARCHAR(255) PRIMARY KEY,
					subdomain VARCHAR(255) NOT NULL,
					domain VARCHAR(255) NOT NULL,
					created_at TIMESTAMP NOT NULL,
					updated_at TIMESTAMP NOT NULL,
					version INT NOT NULL
			)`,
		Down: `DROP TABLE IF EXISTS sites`,
	},
	{
		ID:   "add_user_column",
		Up:   `ALTER TABLE sites ADD COLUMN "user" VARCHAR(255) NOT NULL`,
		Down: `ALTER TABLE sites DROP COLUMN "user"`,
	},
}
