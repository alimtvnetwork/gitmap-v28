package repodb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// InitRepoSchema initializes the repository-specific SQLite DB tables.
func InitRepoSchema(ctx context.Context, db *sql.DB) error {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS RepoFile ( RepoFileId INTEGER PRIMARY KEY AUTOINCREMENT, RelativePath TEXT NOT NULL UNIQUE, AbsolutePath TEXT NOT NULL, Content TEXT, IsBig INTEGER NOT NULL, WriteTime INTEGER NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS SearchCache ( SearchCacheId INTEGER PRIMARY KEY AUTOINCREMENT, Query TEXT NOT NULL UNIQUE, Hits INTEGER NOT NULL, ResultJson TEXT NOT NULL, CreatedAt INTEGER NOT NULL, UpdatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS FileSequence ( FileSequenceId INTEGER PRIMARY KEY AUTOINCREMENT, Directory TEXT NOT NULL, Filename TEXT NOT NULL, SequenceNumber INTEGER NOT NULL, BaseName TEXT NOT NULL, UpdatedAt INTEGER NOT NULL, UNIQUE(Directory, Filename) );",
		"CREATE TABLE IF NOT EXISTS SequenceHistory ( SequenceHistoryId INTEGER PRIMARY KEY AUTOINCREMENT, Directory TEXT NOT NULL, OperationsJson TEXT NOT NULL, CreatedAt INTEGER NOT NULL );",
		"CREATE TABLE IF NOT EXISTS RepoScanLog ( RepoScanLogId INTEGER PRIMARY KEY AUTOINCREMENT, RepoId INTEGER NOT NULL, RepoSlug TEXT NOT NULL, Action TEXT NOT NULL, Status TEXT NOT NULL, ErrorMessage TEXT, Details TEXT, Notes TEXT, Comments TEXT, CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP );",
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return apperror.WrapSimple(err, "init repo schema")
		}
	}

	return nil
}

// ResolveRepoDBPath resolves the target database filepath.
func ResolveRepoDBPath(rootDbDir, absolutePath string, repoId int64) string {
	slug := GenerateSlug(absolutePath)
	repoSearchDir := filepath.Join(rootDbDir, "repo_search")
	if err := os.MkdirAll(repoSearchDir, 0755); err != nil {
		return filepath.Join(rootDbDir, fmt.Sprintf("%s-%d.db", slug, repoId))
	}
	return filepath.Join(repoSearchDir, fmt.Sprintf("%s-%d.db", slug, repoId))
}

func closeAndWrapInitError(db *sql.DB, initErr error) *apperror.AppError {
	if closeErr := db.Close(); closeErr != nil {
		return apperror.WrapWithDetails(
			closeErr,
			"close db after init failure",
			"E9000",
			"close db failed after: "+initErr.Error(),
			"repodb",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
	}
	return apperror.WrapSimple(initErr, "init repo db schema")
}

// OpenRepoDB opens or creates the split DB for a specific repository.
func OpenRepoDB(ctx context.Context, rootDbDir, absolutePath string, repoId int64) (*sql.DB, error) {
	dbPath := ResolveRepoDBPath(rootDbDir, absolutePath, repoId)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open repo db")
	}

	db.SetMaxOpenConns(1)

	if err := InitRepoSchema(ctx, db); err != nil {
		return nil, closeAndWrapInitError(db, err)
	}

	if _, err := dbengine.WrapDb(db, dbengine.DbSQLite); err != nil {
		return nil, closeAndWrapInitError(db, err)
	}

	return db, nil
}

// OpenRepoDbWrapper opens or creates the split DB returning a typed DbWrapper.
func OpenRepoDbWrapper(ctx context.Context, rootDbDir, absolutePath string, repoId int64) (*dbengine.DbWrapper, error) {
	db, err := OpenRepoDB(ctx, rootDbDir, absolutePath, repoId)
	if err != nil {
		return nil, err
	}
	wrap, wrapErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if wrapErr != nil {
		return nil, wrapErr
	}
	return wrap, nil
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
			return apperror.WrapSimple(err, "clear repo db")
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
			return apperror.WrapSimple(err, "reset repo db")
		}
	}

	return InitRepoSchema(ctx, db)
}

func getRepoDBFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func runRepoOptimizePragmas(ctx context.Context, db *sql.DB) *apperror.AppError {
	if _, err := db.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
		return apperror.WrapSimple(err, "wal checkpoint repo db")
	}
	if _, err := db.ExecContext(ctx, "VACUUM;"); err != nil {
		return apperror.WrapSimple(err, "vacuum repo db")
	}
	if _, err := db.ExecContext(ctx, "PRAGMA optimize;"); err != nil {
		return apperror.WrapSimple(err, "optimize repo db")
	}
	return nil
}

// OptimizeRepoDB runs VACUUM and PRAGMA optimize, returning bytes reclaimed.
func OptimizeRepoDB(ctx context.Context, db *sql.DB, path string) (int64, error) {
	sizeBefore := getRepoDBFileSize(path)
	if err := runRepoOptimizePragmas(ctx, db); err != nil {
		return 0, err
	}
	sizeAfter := getRepoDBFileSize(path)
	if sizeBefore <= sizeAfter {
		return 0, nil
	}
	return sizeBefore - sizeAfter, nil
}
