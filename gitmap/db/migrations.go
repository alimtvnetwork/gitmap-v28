package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

// ApplyMigrations applies all SQL migration files in the migrations directory
// against the provided standard database/sql connection.
func ApplyMigrations(ctx context.Context, conn *sql.DB) error {
	entries, err := MigrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || len(entry.Name()) < 4 || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := MigrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
		}

		_, err = conn.ExecContext(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", entry.Name(), err)
		}
	}
	return nil
}
