package uploadsession_test

import (
	"hostamat/internal/uploadsession"
	"slices"
	"testing"
	"time"

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

func TestGet(t *testing.T) {
	db := testDB()
	repo := uploadsession.NewRepository(db)
	repo.Migrate()

	t.Run("non-existent session", func(t *testing.T) {
		_, err := repo.Get("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent session")
		}
	})

	t.Run("existing session", func(t *testing.T) {
		// Create test session
		testSession := &uploadsession.UploadSession{
			Id:        "test-id",
			CreatedAt: time.Now(),
			Status:    uploadsession.Open,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		err := repo.Create(testSession)
		if err != nil {
			t.Fatalf("Failed to create test session: %v", err)
		}

		// Get and verify
		got, err := repo.Get("test-id")
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		if got.Id != testSession.Id {
			t.Errorf("Expected ID %q, got %q", testSession.Id, got.Id)
		}
		if got.Status != testSession.Status {
			t.Errorf("Expected status %q, got %q", testSession.Status, got.Status)
		}
	})
}

func testDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	return db
}
