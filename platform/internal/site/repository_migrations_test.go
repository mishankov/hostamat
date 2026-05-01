package site_test

import (
	"slices"
	"testing"

	"hostamat/internal/site"

	"github.com/platforma-dev/platforma/database"
)

func TestRepositoryMigrations(t *testing.T) {
	t.Parallel()

	repo := site.NewRepository(nil)
	migrations, err := database.ParseMigrations(repo.Migrations())
	if err != nil {
		t.Fatalf("ParseMigrations() error = %v", err)
	}

	ids := make([]string, 0, len(migrations))
	for _, migration := range migrations {
		ids = append(ids, migration.ID)
	}

	expected := []string{"init", "add_user_column", "add_name_column", "drop_version_column", "add_active_release_id"}
	if !slices.Equal(ids, expected) {
		t.Fatalf("migration ids = %v, want %v", ids, expected)
	}
}
