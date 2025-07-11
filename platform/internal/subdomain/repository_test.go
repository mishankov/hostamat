package subdomain_test

import (
	"context"
	"database/sql"
	"errors"
	"hostamat/internal/subdomain"
	"slices"
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
	repo := subdomain.NewRepository(db)

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

		expectedColumns := []string{"id", "name", "full_domain"}

		t.Run("column count", func(t *testing.T) {
			if len(columns) != len(expectedColumns) {
				t.Fatalf("Expected %d columns, got %d", len(expectedColumns), len(columns))
			}
		})

		t.Run("column names", func(t *testing.T) {
			for i, column := range columns {
				if !slices.Contains(expectedColumns, column) {
					t.Errorf("Expected column %s to be in %s", column, expectedColumns[i])
				}
			}
		})
	})
}

func TestRepository_Create(t *testing.T) {
	db := testDB(t)
	repo := subdomain.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Create(context.Background(), &subdomain.Subdomain{
			Id:   "test-id",
			Name: "test-name",
		}); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		var count int
		if err := db.Get(&count, "SELECT COUNT(*) FROM subdomains WHERE id = ?", "test-id"); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 record, got %d", count)
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		err := repo.Create(context.Background(), &subdomain.Subdomain{
			Id:   "test-id-2",
			Name: "test-name2",
		})
		if err != nil {
			t.Fatal(err)
		}

		err = repo.Create(context.Background(), &subdomain.Subdomain{
			Id:   "test-id-2",
			Name: "test-name22",
		})
		if err == nil {
			t.Error("expected error for duplicate id")
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Create fresh subdomain with unique id to avoid id conflict
		if err := repo.Create(context.Background(), &subdomain.Subdomain{
			Id:   "test-id3",
			Name: "test-name3",
		}); err != nil {
			t.Fatal(err)
		}

		err := repo.Create(context.Background(), &subdomain.Subdomain{
			Id:   "test-id4",
			Name: "test-name3",
		})
		if err == nil {
			t.Error("expected error for duplicate name")
		}
	})
}

func TestRepository_Get(t *testing.T) {
	db := testDB(t)
	repo := subdomain.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSubdomain := &subdomain.Subdomain{
		Id:   "test-id",
		Name: "test-name",
	}

	t.Run("not found", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "nonexistent")
		if !errors.Is(err, sql.ErrNoRows) {
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
