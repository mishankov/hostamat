package subdomain

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func testDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func TestRepository_Migrate(t *testing.T) {
	db := testDB()
	repo := NewRepository(db)

	err := repo.Migrate()
	assert.NoError(t, err)

	t.Run("table exists", func(t *testing.T) {
		rows, err := db.Query("SELECT * FROM subdomains")
		assert.NoError(t, err)
		defer rows.Close()

		columns, err := rows.Columns()
		assert.NoError(t, err)
		assert.Equal(t, []string{"id", "name"}, columns)
	})
}

func TestRepository_Create(t *testing.T) {
	db := testDB()
	repo := NewRepository(db)
	repo.Migrate()

	testSubdomain := &Subdomain{
		Id:   "test-id",
		Name: "test-name",
	}

	t.Run("success", func(t *testing.T) {
		err := repo.Create(context.Background(), testSubdomain)
		assert.NoError(t, err)

		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM subdomains WHERE id = ?", testSubdomain.Id)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("duplicate id", func(t *testing.T) {
		err := repo.Create(context.Background(), testSubdomain)
		assert.NoError(t, err)

		err = repo.Create(context.Background(), testSubdomain)
		assert.Error(t, err)
	})
}

func TestRepository_Get(t *testing.T) {
	db := testDB()
	repo := NewRepository(db)
	repo.Migrate()

	testSubdomain := &Subdomain{
		Id:   "test-id",
		Name: "test-name",
	}

	t.Run("not found", func(t *testing.T) {
		_, err := repo.Get(context.Background(), "nonexistent")
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("found", func(t *testing.T) {
		err := repo.Create(context.Background(), testSubdomain)
		assert.NoError(t, err)

		got, err := repo.Get(context.Background(), testSubdomain.Id)
		assert.NoError(t, err)
		assert.Equal(t, testSubdomain.Id, got.Id)
		assert.Equal(t, testSubdomain.Name, got.Name)
	})
}
