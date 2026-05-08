package formcollection_test

import (
	"slices"
	"testing"

	"hostamat/internal/formcollection"

	"github.com/platforma-dev/platforma/database"
)

func TestRepositoryMigrations(t *testing.T) {
	t.Parallel()

	repo := formcollection.NewRepository(nil)
	migrations, err := database.ParseMigrations(repo.Migrations())
	if err != nil {
		t.Fatalf("ParseMigrations() error = %v", err)
	}

	ids := make([]string, 0, len(migrations))
	for _, migration := range migrations {
		ids = append(ids, migration.ID)
	}

	expected := []string{"init"}
	if !slices.Equal(ids, expected) {
		t.Fatalf("migration ids = %v, want %v", ids, expected)
	}
}
