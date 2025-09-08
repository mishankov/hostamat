package database

import "time"

type migrations struct {
	Repository  string    `db:"repository"`
	MigrationId string    `db:"id"`
	Timestamp   time.Time `db:"timestamp"`
}

type Migration struct {
	ID         string
	Up         string
	Down       string
	repository string
}

type Schema struct {
	Queries []string
}

type shemer interface {
	Schema() ([]Migration, Schema)
}
