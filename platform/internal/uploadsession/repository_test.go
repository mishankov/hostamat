package uploadsession_test

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"
	"time"

	"hostamat/internal/uploadsession"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
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
	repo := uploadsession.NewRepository(db)

	if err := repo.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	t.Run("table exists", func(t *testing.T) {
		rows, err := db.Query("SELECT * FROM upload_sessions")
		if err != nil {
			t.Fatalf("failed to query table: %v", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			t.Fatalf("failed to get columns: %v", err)
		}

		expectedColumns := []string{"id", "site_id", "status", "created_at", "updated_at"}

		t.Run("column count", func(t *testing.T) {
			if len(columns) != len(expectedColumns) {
				t.Fatalf("Expected %d columns, got %d", len(expectedColumns), len(columns))
			}
		})

		t.Run("column names", func(t *testing.T) {
			for _, column := range columns {
				if !slices.Contains(expectedColumns, column) {
					t.Errorf("Unexpected column %s", column)
				}
			}
		})
	})
}

func TestRepository_Create(t *testing.T) {
	db := testDB(t)
	repo := uploadsession.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSession := &uploadsession.UploadSession{
		Id:     uuid.NewString(),
		SiteId: "test-site-id",
		Status: uploadsession.StatusOpen,
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSession); err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := repo.Get(context.Background(), testSession.Id)
		if err != nil {
			t.Fatal(err)
		}
		if got.Id != testSession.Id {
			t.Errorf("expected id %q, got %q", testSession.Id, got.Id)
		}
	})
}

func TestRepository_Get(t *testing.T) {
	db := testDB(t)
	repo := uploadsession.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSession := &uploadsession.UploadSession{
		Id:     uuid.NewString(),
		SiteId: "test-site-id",
		Status: uploadsession.StatusOpen,
	}

	t.Run("not found", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "nonexistent")
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("found", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSession); err != nil {
			t.Fatal(err)
		}

		got, err := repo.Get(context.Background(), testSession.Id)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.Id != testSession.Id {
			t.Errorf("expected id %q, got %q", testSession.Id, got.Id)
		}
		if got.SiteId != testSession.SiteId {
			t.Errorf("expected siteId %q, got %q", testSession.SiteId, got.SiteId)
		}
		if got.Status != testSession.Status {
			t.Errorf("expected status %q, got %q", testSession.Status, got.Status)
		}
		if got.CreatedAt.IsZero() {
			t.Error("expected createdAt to be set")
		}
		if got.UpdatedAt.IsZero() {
			t.Error("expected updatedAt to be set")
		}
	})
}

func TestRepository_UpdateStatus(t *testing.T) {
	db := testDB(t)
	repo := uploadsession.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		t.Fatal(err)
	}

	testSession := &uploadsession.UploadSession{
		Id:     uuid.NewString(),
		SiteId: "test-site-id",
		Status: uploadsession.StatusOpen,
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Create(context.Background(), testSession); err != nil {
			t.Fatal(err)
		}

		// Get initial timestamps
		initial, err := repo.Get(context.Background(), testSession.Id)
		if err != nil {
			t.Fatal(err)
		}

		// Wait briefly to ensure timestamp would change
		time.Sleep(10 * time.Millisecond)

		if err := repo.UpdateStatus(context.Background(), testSession.Id, uploadsession.StatusClosed); err != nil {
			t.Fatalf("UpdateStatus() error = %v", err)
		}

		updated, err := repo.Get(context.Background(), testSession.Id)
		if err != nil {
			t.Fatal(err)
		}

		if updated.Status != uploadsession.StatusClosed {
			t.Errorf("expected status %q, got %q", uploadsession.StatusClosed, updated.Status)
		}
		if !updated.UpdatedAt.After(initial.UpdatedAt) {
			t.Errorf("expected updatedAt to be updated, initial: %v, updated: %v", initial.UpdatedAt, updated.UpdatedAt)
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.UpdateStatus(context.Background(), "nonexistent", uploadsession.StatusClosed)
		if err == nil {
			t.Error("expected error for nonexistent session")
		}
	})
}
