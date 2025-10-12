package site_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"hostamat/internal/site"
	"hostamat/pkg/platforma/database"

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

func setupTestDB(t *testing.T) (*sqlx.DB, *site.Repository) {
	db := testDB(t)
	repo := site.NewRepository(db)

	// Setup database with migrations
	dbObj, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	dbObj.DB = db // Use our test db
	dbObj.RegisterRepository("sites", repo)

	if err := dbObj.Migrate(context.Background()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db, repo
}

func TestRepository_Create(t *testing.T) {
	db, repo := setupTestDB(t)

	t.Run("successful creation", func(t *testing.T) {
		testSite := &site.Site{
			Id:        "create-test-id",
			User:      "test-user",
			Subdomain: "example",
			Domain:    "example.example.com",
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
			User:      "test-user",
			Subdomain: "example",
			Domain:    "example.example.com",
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

func TestRepository_Get(t *testing.T) {
	_, repo := setupTestDB(t)

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
			User:      "test-user",
			Subdomain: "example",
			Domain:    "example.example.com",
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
		if got.User != testSite.User {
			t.Errorf("Expected user %q, got %q", testSite.User, got.User)
		}
		if got.Subdomain != testSite.Subdomain {
			t.Errorf("Expected subdomain %q, got %q", testSite.Subdomain, got.Subdomain)
		}
		if got.Domain != testSite.Domain {
			t.Errorf("Expected domain %q, got %q", testSite.Domain, got.Domain)
		}
		if !got.CreatedAt.Equal(testSite.CreatedAt) {
			t.Errorf("Expected CreatedAt %v, got %v", testSite.CreatedAt, got.CreatedAt)
		}
	})
}

func TestRepository_GetList(t *testing.T) {
	_, repo := setupTestDB(t)

	// Create test sites for different users
	testSites := []*site.Site{
		{
			Id:        "1",
			User:      "user1",
			Subdomain: "site1",
			Domain:    "site1.example.com",
			CreatedAt: time.Now(),
		},
		{
			Id:        "2",
			User:      "user1",
			Subdomain: "site2",
			Domain:    "site2.example.com",
			CreatedAt: time.Now(),
		},
		{
			Id:        "3",
			User:      "user2",
			Subdomain: "site3",
			Domain:    "site3.example.com",
			CreatedAt: time.Now(),
		},
	}

	for _, s := range testSites {
		if err := repo.Create(context.Background(), s); err != nil {
			t.Fatalf("Failed to create test site: %v", err)
		}
	}

	t.Run("get list for user1", func(t *testing.T) {
		// Test GetList for user1
		sites, err := repo.GetList(context.Background(), "user1")
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if len(sites) != 2 {
			t.Fatalf("Expected 2 sites for user1, got %d", len(sites))
		}

		// Verify only user1's sites are returned
		for _, s := range sites {
			if s.User != "user1" {
				t.Errorf("Expected only user1 sites, got user %q", s.User)
			}
		}
	})

	t.Run("get list for user2", func(t *testing.T) {
		// Test GetList for user2
		sites, err := repo.GetList(context.Background(), "user2")
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if len(sites) != 1 {
			t.Fatalf("Expected 1 site for user2, got %d", len(sites))
		}

		if sites[0].User != "user2" {
			t.Errorf("Expected user2, got user %q", sites[0].User)
		}
	})

	t.Run("get list for non-existent user", func(t *testing.T) {
		// Test GetList for user with no sites
		sites, err := repo.GetList(context.Background(), "nonexistent-user")
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if len(sites) != 0 {
			t.Fatalf("Expected 0 sites for nonexistent user, got %d", len(sites))
		}
	})
}

func TestRepository_GetByDomain(t *testing.T) {
	_, repo := setupTestDB(t)

	t.Run("non-existent domain", func(t *testing.T) {
		_, err := repo.GetByDomain(context.Background(), "nonexistent.example.com")
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Expected sql.ErrNoRows or nil, got %v", err)
		}
	})

	t.Run("existing domain", func(t *testing.T) {
		// Create test site
		testSite := &site.Site{
			Id:        "test-id",
			User:      "test-user",
			Subdomain: "example",
			Domain:    "example.example.com",
			CreatedAt: time.Now(),
		}
		err := repo.Create(context.Background(), testSite)
		if err != nil {
			t.Fatalf("Failed to create test site: %v", err)
		}

		// Get by domain and verify
		got, err := repo.GetByDomain(context.Background(), "example.example.com")
		if err != nil {
			t.Fatalf("Failed to get site by domain: %v", err)
		}

		if got.Id != testSite.Id {
			t.Errorf("Expected ID %q, got %q", testSite.Id, got.Id)
		}
		if got.Domain != testSite.Domain {
			t.Errorf("Expected domain %q, got %q", testSite.Domain, got.Domain)
		}
	})
}
