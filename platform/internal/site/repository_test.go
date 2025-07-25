package site_test

import (
	"context"
	"testing"
	"time"

	"hostamat/internal/site"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func TestMigrate(t *testing.T) {
	db := testDB()
	repo := site.NewRepository(db)
	repo.Migrate()

	t.Run("table exists", func(t *testing.T) {
		rows, err := db.Query("SELECT * FROM sites")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}

		expectedColumns := []string{"id", "subdomain", "domain", "created_at", "updated_at", "version"}

		if len(columns) != len(expectedColumns) {
			t.Fatalf("Expected %d columns, got %d", len(expectedColumns), len(columns))
		}

		for i, column := range columns {
			if column != expectedColumns[i] {
				t.Errorf("Expected column %s, got %s", expectedColumns[i], column)
			}
		}
	})
}

func TestCreate(t *testing.T) {
	db := testDB()
	repo := site.NewRepository(db)
	repo.Migrate()

	t.Run("successful creation", func(t *testing.T) {
		testSite := &site.Site{
			Id:        "create-test-id",
			Subdomain: "example",
			CreatedAt: time.Now(),
		}

		err := repo.Create(context.Background(), testSite)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Verify the site was actually created
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM sites WHERE id = ?", testSite.Id)
		if err != nil {
			t.Fatalf("Failed to verify creation: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 record, got %d", count)
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		testSite := &site.Site{
			Id:        "duplicate-id",
			Subdomain: "example",
			CreatedAt: time.Now(),
		}

		// First create should succeed
		err := repo.Create(context.Background(), testSite)
		if err != nil {
			t.Fatalf("First create failed: %v", err)
		}

		// Second create with same ID should fail
		err = repo.Create(context.Background(), testSite)
		if err == nil {
			t.Error("Expected error for duplicate ID")
		}
	})
}

func TestGet(t *testing.T) {
	db := testDB()
	repo := site.NewRepository(db)
	repo.Migrate()

	t.Run("non-existent site", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent site")
		}
	})

	t.Run("existing site", func(t *testing.T) {
		// Create test site
		testSite := &site.Site{
			Id:        "test-id",
			Subdomain: "example",
			CreatedAt: time.Now(),
		}
		err := repo.Create(context.Background(), testSite)
		if err != nil {
			t.Fatalf("Failed to create test site: %v", err)
		}

		// Get and verify
		got, err := repo.Get(context.Background(), "test-id")
		if err != nil {
			t.Fatalf("Failed to get site: %v", err)
		}

		// Compare all fields
		if got.Id != testSite.Id {
			t.Errorf("Expected ID %q, got %q", testSite.Id, got.Id)
		}
		if got.Subdomain != testSite.Subdomain {
			t.Errorf("Expected subdomain %q, got %q", testSite.Subdomain, got.Subdomain)
		}
		if !got.CreatedAt.Equal(testSite.CreatedAt) {
			t.Errorf("Expected CreatedAt %v, got %v", testSite.CreatedAt, got.CreatedAt)
		}
	})
}

func TestGetList(t *testing.T) {
	db := testDB()
	repo := site.NewRepository(db)
	repo.Migrate()

	// Create test sites
	testSites := []*site.Site{
		{
			Id:        "1",
			Subdomain: "site1",
			CreatedAt: time.Now(),
		},
		{
			Id:        "2",
			Subdomain: "site2",
			CreatedAt: time.Now(),
		},
	}

	for _, s := range testSites {
		if err := repo.Create(context.Background(), s); err != nil {
			t.Fatalf("Failed to create test site: %v", err)
		}
	}

	// Test GetList
	sites, err := repo.GetList(context.Background())
	if err != nil {
		t.Fatalf("GetList failed: %v", err)
	}

	if len(sites) != len(testSites) {
		t.Fatalf("Expected %d sites, got %d", len(testSites), len(sites))
	}

	for i, s := range sites {
		if s.Id != testSites[i].Id {
			t.Errorf("Expected ID %q, got %q", testSites[i].Id, s.Id)
		}
		if s.Subdomain != testSites[i].Subdomain {
			t.Errorf("Expected subdomain %q, got %q", testSites[i].Subdomain, s.Subdomain)
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
