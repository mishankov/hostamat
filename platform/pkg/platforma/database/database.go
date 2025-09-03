package database

import (
	"context"
	"hostamat/pkg/platforma/log"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
	repositories map[string]any
	migrators    map[string]shemer
}

func New(driver, connection string) (*Database, error) {
	db, err := sqlx.Connect(driver, connection)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db, repositories: make(map[string]any), migrators: make(map[string]shemer)}, nil
}

func (db *Database) RegisterRepository(name string, repository any) {
	db.repositories[name] = repository

	if migr, ok := repository.(shemer); ok {
		db.migrators[name] = migr
	}
}

func (db *Database) Migrate(ctx context.Context) error {
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS platforma_migrations (repository TEXT, id TEXT, timestamp DATETIME)"); err != nil {
		return err
	}

	// Select data from platforma_migrations table
	var migrationsState []migrations
	err := db.SelectContext(ctx, &migrationsState, "SELECT * FROM platforma_migrations")
	if err != nil {
		return err
	}

	for repoName, migr := range db.migrators {
		repoMigrations, repoSchema := migr.Schema()

		repoHasMigrations := slices.ContainsFunc(migrationsState, func(m migrations) bool {
			return m.Repository == repoName
		})

		// If repo does not has migrations apply schema and exit
		if !repoHasMigrations {
			for _, query := range repoSchema.Queries {
				if _, err := db.ExecContext(ctx, query); err != nil {
					return err
				}
				log.InfoContext(ctx, "schema applied", "repository", repoName)
			}

			// Log that schema applied
			if _, err := db.ExecContext(ctx, "INSERT INTO platforma_migrations (repository, timestamp) VALUES (?, ?)", repoName, time.Now()); err != nil {
				return err
			}

			continue
		}

		for _, migration := range repoMigrations {
			// Check if migration has been applied
			migrationHasApplied := slices.ContainsFunc(migrationsState, func(m migrations) bool {
				return m.Repository == repoName && m.MigrationId == migration.ID
			})

			if migrationHasApplied {
				continue
			}

			for _, query := range migration.Queries {
				if _, err := db.ExecContext(ctx, query); err != nil {
					return err
				}

				log.InfoContext(ctx, "applied migration %s for repository %s", migration.ID, repoName)
			}

			// Log that migration applied
			if _, err := db.ExecContext(ctx, "INSERT INTO platforma_migrations (repository, migration_id, timestamp) VALUES (?, ?, ?)", repoName, migration.ID, time.Now()); err != nil {
				return err
			}
		}
	}

	return nil
}
