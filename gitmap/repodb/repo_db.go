package repodb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// InitRepoSchema initializes the repository-specific SQLite DB tables

func InitRepoSchema(ctx context.Context, db *sql.DB) error {
	query1 := "CREATE TABLE IF NOT EXISTS RepoFile ( Id INTEGER PRIMARY KEY AUTOINCREMENT, RelativePath TEXT NOT NULL UNIQUE, AbsolutePath TEXT NOT NULL, Content TEXT, IsBig INTEGER NOT NULL, WriteTime INTEGER NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );"
	query2 := "CREATE TABLE IF NOT EXISTS SearchCache ( Id INTEGER PRIMARY KEY AUTOINCREMENT, Query TEXT NOT NULL UNIQUE, Hits INTEGER NOT NULL, ResultJson TEXT NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );"

	if _, err := db.ExecContext(ctx, query1); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, query2)
	return err
}

// OpenRepoDB opens or creates the split DB for a specific repository.

func OpenRepoDB(
	ctx context.Context,
	rootDbDir,
	absolutePath string,
	repoId int64,
) (*sql.DB, error) {
	slug := GenerateSlug(absolutePath)
	repoSearchDir := filepath.Join(rootDbDir, "repo_search")
	if err := os.MkdirAll(repoSearchDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(repoSearchDir, fmt.Sprintf("%s-%d.db", slug, repoId))

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := InitRepoSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
