package database

import "time"

type migrations struct {
	Repository  string    `db:"repository"`
	MigrationId string    `db:"id"`
	Timestamp   time.Time `db:"timestamp"`
}

type Migration struct {
	ID      string
	Queries []string
}

type Schema struct {
	Queries []string
}

type migrator interface {
	Schema() ([]Migration, Schema)
}
