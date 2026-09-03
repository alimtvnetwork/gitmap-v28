package repodb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// InitRepoSchema initializes the repository-specific SQLite DB tables.
func InitRepoSchema(ctx context.Context, db *sql.DB) error {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS RepoFile ( Id INTEGER PRIMARY KEY AUTOINCREMENT, RelativePath TEXT NOT NULL UNIQUE, AbsolutePath TEXT NOT NULL, Content TEXT, IsBig INTEGER NOT NULL, WriteTime INTEGER NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS SearchCache ( Id INTEGER PRIMARY KEY AUTOINCREMENT, Query TEXT NOT NULL UNIQUE, Hits INTEGER NOT NULL, ResultJson TEXT NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS FileSequence ( Id INTEGER PRIMARY KEY AUTOINCREMENT, Directory TEXT NOT NULL, Filename TEXT NOT NULL, SequenceNumber INTEGER NOT NULL, BaseName TEXT NOT NULL, UpdatedAt INTEGER NOT NULL, UNIQUE(Directory, Filename) );",
		"CREATE TABLE IF NOT EXISTS SequenceHistory ( Id INTEGER PRIMARY KEY AUTOINCREMENT, Directory TEXT NOT NULL, OperationsJson TEXT NOT NULL, CreatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS RepoScanLog ( Id INTEGER PRIMARY KEY AUTOINCREMENT, RepoId INTEGER NOT NULL, RepoSlug TEXT NOT NULL, Action TEXT NOT NULL, Status TEXT NOT NULL, ErrorMessage TEXT, Details TEXT, Notes TEXT, Comments TEXT, CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP );",
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

// ResolveRepoDBPath resolves the target database filepath.
func ResolveRepoDBPath(rootDbDir, absolutePath string, repoId int64) string {
	slug := GenerateSlug(absolutePath)
	repoSearchDir := filepath.Join(rootDbDir, "repo_search")
	_ = os.MkdirAll(repoSearchDir, 0755)
	return filepath.Join(repoSearchDir, fmt.Sprintf("%s-%d.db", slug, repoId))
}

// OpenRepoDB opens or creates the split DB for a specific repository.
func OpenRepoDB(ctx context.Context, rootDbDir, absolutePath string, repoId int64) (*sql.DB, error) {
	dbPath := ResolveRepoDBPath(rootDbDir, absolutePath, repoId)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if err := InitRepoSchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// ClearRepoDB clears cached queries and search indexes.
func ClearRepoDB(ctx context.Context, db *sql.DB) error {
	queries := []string{
		"DELETE FROM SearchCache;",
		"DELETE FROM RepoFile;",
		"DELETE FROM FileSequence;",
	}
	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

// ResetRepoDB drops all tables and rebuilds fresh schema.
func ResetRepoDB(ctx context.Context, db *sql.DB) error {
	queries := []string{
		"DROP TABLE IF EXISTS RepoFile;",
		"DROP TABLE IF EXISTS SearchCache;",
		"DROP TABLE IF EXISTS FileSequence;",
		"DROP TABLE IF EXISTS SequenceHistory;",
		"DROP TABLE IF EXISTS RepoScanLog;",
	}
	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return InitRepoSchema(ctx, db)
}

// OptimizeRepoDB runs VACUUM and PRAGMA optimize, returning bytes reclaimed.
func OptimizeRepoDB(ctx context.Context, db *sql.DB, path string) (int64, error) {
	var sizeBefore int64
	if info, err := os.Stat(path); err == nil {
		sizeBefore = info.Size()
	}
	_, _ = db.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE);")
	if _, err := db.ExecContext(ctx, "VACUUM;"); err != nil {
		return 0, err
	}
	_, _ = db.ExecContext(ctx, "PRAGMA optimize;")
	var sizeAfter int64
	if info, err := os.Stat(path); err == nil {
		sizeAfter = info.Size()
	}
	reclaimed := sizeBefore - sizeAfter
	if reclaimed < 0 {
		reclaimed = 0
	}
	return reclaimed, nil
}
