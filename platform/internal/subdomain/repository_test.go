package subdomain

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func testDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	return db
}

func TestRepository_Migrate(t *testing.T) {
	db := testDB(t)
	repo := New(db)

	if err := repo.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	t.Run("table exists", func(t *testing.T) {
		rows, err := db.Query("SELECT * FROM subdomains")
		if err != nil {
			t.Fatalf("failed to query table: %v", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			t.Fatalf("failed to get columns: %v", err)
		}

		if len(columns) != 2 {
			t.Fatalf("expected 2 columns, got %d", len(columns))
		}
		if columns[0] != "id" || columns[1] != "name" {
			t.Fatalf("expected columns [id name], got %v", columns)
		}
	})
}

func TestRepository_Create(t *testing.T) {
	db := testDB(t)
	repo := New(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSubdomain := &Subdomain{
		Id:   "test-id",
		Name: "test-name",
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSubdomain); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		var count int
		if err := db.Get(&count, "SELECT COUNT(*) FROM subdomains WHERE id = ?", testSubdomain.Id); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 record, got %d", count)
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSubdomain); err != nil {
			t.Fatal(err)
		}

		err := repo.Create(context.Background(), testSubdomain)
		if err == nil {
			t.Error("expected error for duplicate id")
		} else if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			t.Errorf("expected constraint error, got %v", err)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	db := testDB(t)
	repo := New(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSubdomain := &Subdomain{
		Id:   "test-id",
		Name: "test-name",
	}

	t.Run("not found", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "nonexistent")
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("found", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSubdomain); err != nil {
			t.Fatal(err)
		}

		got, err := repo.Get(context.Background(), testSubdomain.Id)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.Id != testSubdomain.Id {
			t.Errorf("expected id %q, got %q", testSubdomain.Id, got.Id)
		}
		if got.Name != testSubdomain.Name {
			t.Errorf("expected name %q, got %q", testSubdomain.Name, got.Name)
		}
	})
}
