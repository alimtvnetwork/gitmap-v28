package db

import (
	"context"
	"database/sql"
	"embed"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

// ApplyMigrations applies all SQL migration files in the migrations directory
// against the provided standard database/sql connection.

func ApplyMigrations(ctx context.Context, conn *sql.DB) error {
	entries, err := MigrationsFS.ReadDir("migrations")
	if err != nil {
		return apperror.WrapSimple(err, "failed to read migrations dir")
	}

	for _, entry := range entries {
		if entry.IsDir() || len(entry.Name()) < 4 || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := MigrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return apperror.WrapSimple(err, "failed to read migration "+entry.Name())
		}

		_, err = conn.ExecContext(ctx, string(content))
		if err != nil {
			return apperror.WrapSimple(err, "failed to execute migration "+entry.Name())
		}
	}

	return nil
}
