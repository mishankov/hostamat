package customdomain_test

import (
	"slices"
	"testing"

	"hostamat/internal/customdomain"

	"github.com/platforma-dev/platforma/database"
)

func TestRepositoryMigrations(t *testing.T) {
	t.Parallel()

	repo := customdomain.NewRepository(nil)
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
