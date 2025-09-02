package migrations

import "hostamat/pkg/platforma/database"

var Migrations = []database.Migration{
	{
		ID: "init",
		Queries: []string{
			`CREATE TABLE IF NOT EXISTS sites (
				id VARCHAR(255) PRIMARY KEY,
				subdomain VARCHAR(255) NOT NULL,
				domain VARCHAR(255) NOT NULL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				version INT NOT NULL
			)`,
		},
	},
}

var Schema = database.Schema{
	Queries: []string{
		`CREATE TABLE IF NOT EXISTS sites (
			id VARCHAR(255) PRIMARY KEY,
			subdomain VARCHAR(255) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			version INT NOT NULL
		)`,
	},
}
