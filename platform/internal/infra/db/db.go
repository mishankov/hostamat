package db

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func CreateDB(connection string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}
