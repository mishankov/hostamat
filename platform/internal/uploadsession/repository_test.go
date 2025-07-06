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

	t.Run("table exists", func(t *testing.T) {
		rows, err := db.Query("SELECT * FROM upload_sessions")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		t.Run("has correct columns", func(t *testing.T) {
			columns, err := rows.Columns()
			if err != nil {
				t.Fatal(err)
			}

			expectedColumns := []string{"id", "created_at", "expires_at", "status", "resolved_at"}

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
	})
}

func TestCreate(t *testing.T) {
	db := testDB()
	repo := uploadsession.NewRepository(db)
	repo.Migrate()

	t.Run("successful creation", func(t *testing.T) {
		testSession := &uploadsession.UploadSession{
			Id:        "create-test-id",
			CreatedAt: time.Now(),
			Status:    uploadsession.Open,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		err := repo.Create(testSession)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Verify the session was actually created
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM upload_sessions WHERE id = ?", testSession.Id)
		if err != nil {
			t.Fatalf("Failed to verify creation: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 record, got %d", count)
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		testSession := &uploadsession.UploadSession{
			Id:        "duplicate-id",
			CreatedAt: time.Now(),
			Status:    uploadsession.Open,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		// First create should succeed
		err := repo.Create(testSession)
		if err != nil {
			t.Fatalf("First create failed: %v", err)
		}

		// Second create with same ID should fail
		err = repo.Create(testSession)
		if err == nil {
			t.Error("Expected error for duplicate ID")
		}
	})
}

func TestUpdate(t *testing.T) {
	db := testDB()
	repo := uploadsession.NewRepository(db)
	repo.Migrate()

	t.Run("update status", func(t *testing.T) {
		// Create test session
		testSession := &uploadsession.UploadSession{
			Id:        "update-test-id",
			CreatedAt: time.Now(),
			Status:    uploadsession.Open,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		err := repo.Create(testSession)
		if err != nil {
			t.Fatalf("Failed to create test session: %v", err)
		}

		// Update the session
		updatedSession := *testSession
		updatedSession.Status = uploadsession.Finished
		updatedSession.ResolvedAt = time.Now()

		err = repo.Update(&updatedSession)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Verify the update
		got, err := repo.Get("update-test-id")
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}

		if got.Status != updatedSession.Status {
			t.Errorf("Expected status %q, got %q", updatedSession.Status, got.Status)
		}
		if !got.ResolvedAt.Equal(updatedSession.ResolvedAt) {
			t.Errorf("Expected ResolvedAt %v, got %v", updatedSession.ResolvedAt, got.ResolvedAt)
		}
	})
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

		// Compare all fields
		if got.Id != testSession.Id {
			t.Errorf("Expected ID %q, got %q", testSession.Id, got.Id)
		}
		if !got.CreatedAt.Equal(testSession.CreatedAt) {
			t.Errorf("Expected CreatedAt %v, got %v", testSession.CreatedAt, got.CreatedAt)
		}
		if got.Status != testSession.Status {
			t.Errorf("Expected status %q, got %q", testSession.Status, got.Status)
		}
		if !got.ExpiresAt.Equal(testSession.ExpiresAt) {
			t.Errorf("Expected ExpiresAt %v, got %v", testSession.ExpiresAt, got.ExpiresAt)
		}
		if !got.ResolvedAt.IsZero() {
			t.Errorf("Expected empty ResolvedAt, got %v", got.ResolvedAt)
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
