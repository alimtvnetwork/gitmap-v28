package repodb

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
)

// InitRootSchema initializes the required tables in the root database for SplitDB

func InitRootSchema(ctx context.Context, db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS IndexedRepo ( Id INTEGER PRIMARY KEY AUTOINCREMENT, Path TEXT NOT NULL UNIQUE, Slug TEXT NOT NULL, MigratedVersion INTEGER NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );"
	_, err := db.ExecContext(ctx, query)

	return err
}

// GenerateSlug sanitizes the folder name to be used in DB file names.
func GenerateSlug(absolutePath string) string {
	base := filepath.Base(absolutePath)
	slug := strings.ReplaceAll(base, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "-")

	return slug
}
