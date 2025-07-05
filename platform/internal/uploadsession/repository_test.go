package uploadsession_test

import (
	"hostamat/internal/uploadsession"
	"slices"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func TestMigrate(t *testing.T) {
	db := testDB()
	repo := uploadsession.NewRepository(db)
	repo.Migrate()

	rows, err := db.Query("SELECT * FROM upload_sessions")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}

	expectedColumns := []string{"id", "created_at", "expires_at", "status", "resolved_at"}

	if len(columns) != len(expectedColumns) {
		t.Fatalf("Expected %d columns, got %d", len(expectedColumns), len(columns))
	}

	for i, column := range columns {
		if !slices.Contains(expectedColumns, column) {
			t.Fatalf("Expected column %s to be in %s", column, expectedColumns[i])
		}
	}
}

func testDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	return db
}
