package cmdfoldertree

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	DefaultTreeDBRelativePath = ".gitmap/data/treedb/sql.db"
	FallbackTreeDBPath        = ".gitmap/treedb/tree_cache.db"
)

// ResolveTreeDBPath returns the resolved filesystem path for the tree database.
func ResolveTreeDBPath(customPath string) string {
	if strings.TrimSpace(customPath) != "" {
		return customPath
	}
	// Check primary location
	if _, err := os.Stat(DefaultTreeDBRelativePath); err == nil {
		return DefaultTreeDBRelativePath
	}
	// Check fallback
	if _, err := os.Stat(FallbackTreeDBPath); err == nil {
		return FallbackTreeDBPath
	}
	return DefaultTreeDBRelativePath
}

// OpenTreeDB opens and initializes the Split SQLite database with WAL mode and pragmas.
func OpenTreeDB(dbPath string) (*sql.DB, error) {
	resolved := ResolveTreeDBPath(dbPath)
	dir := filepath.Dir(resolved)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, apperror.WrapSimple(err, "failed to create database directory "+dir)
		}
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", resolved)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to open sqlite database at "+resolved)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	if err := initTreeDBSchema(db); err != nil {
		db.Close()
		return nil, apperror.WrapSimple(err, "failed to initialize schema in "+resolved)
	}

	return db, nil
}

// initTreeDBSchema creates the TreeFile table and case-insensitive indexes.
func initTreeDBSchema(db *sql.DB) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS TreeFile (
		TreeFileId INTEGER PRIMARY KEY AUTOINCREMENT,
		RootPath TEXT,
		RelPath TEXT UNIQUE,
		FileName TEXT,
		DirPath TEXT,
		Extension TEXT,
		SizeBytes INTEGER,
		IsDir INTEGER,
		Depth INTEGER,
		UpdatedAt INTEGER
	);

	CREATE INDEX IF NOT EXISTS IdxTreeFile_RelPath ON TreeFile(RelPath COLLATE NOCASE);
	CREATE INDEX IF NOT EXISTS IdxTreeFile_FileName ON TreeFile(FileName COLLATE NOCASE);
	CREATE INDEX IF NOT EXISTS IdxTreeFile_DirPath ON TreeFile(DirPath COLLATE NOCASE);
	CREATE INDEX IF NOT EXISTS IdxTreeFile_Extension ON TreeFile(Extension COLLATE NOCASE);
	CREATE INDEX IF NOT EXISTS IdxTreeFile_IsDir ON TreeFile(IsDir);

	CREATE VIEW IF NOT EXISTS ViewTreeFiles AS
	SELECT 
		TreeFileId,
		RootPath,
		RelPath,
		FileName,
		DirPath,
		Extension,
		SizeBytes,
		IsDir,
		Depth,
		UpdatedAt
	FROM TreeFile
	ORDER BY DirPath, FileName;
	`

	_, err := db.Exec(ddl)
	return err
}

// IngestTreeItems stores a slice of TreeItem records transactionally into the SQLite database.
func IngestTreeItems(dbPath string, items []TreeItem) (int, error) {
	if len(items) == 0 {
		return 0, nil
	}

	db, err := OpenTreeDB(dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, apperror.WrapSimple(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO TreeFile (
			RootPath, RelPath, FileName, DirPath, Extension, SizeBytes, IsDir, Depth, UpdatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(RelPath) DO UPDATE SET
			RootPath = excluded.RootPath,
			FileName = excluded.FileName,
			DirPath = excluded.DirPath,
			Extension = excluded.Extension,
			SizeBytes = excluded.SizeBytes,
			IsDir = excluded.IsDir,
			Depth = excluded.Depth,
			UpdatedAt = excluded.UpdatedAt;
	`)
	if err != nil {
		return 0, apperror.WrapSimple(err, "failed to prepare upsert statement")
	}
	defer stmt.Close()

	nowUnix := time.Now().Unix()
	ingestedCount := 0

	for _, item := range items {
		isDirInt := 0
		if item.IsDir {
			isDirInt = 1
		}
		updatedAt := item.UpdatedAt
		if updatedAt == 0 {
			updatedAt = nowUnix
		}

		_, err := stmt.Exec(
			item.RootPath,
			item.RelPath,
			item.FileName,
			item.DirPath,
			item.Extension,
			item.SizeBytes,
			isDirInt,
			item.Depth,
			updatedAt,
		)
		if err != nil {
			return ingestedCount, apperror.WrapSimple(err, "failed to upsert item "+item.RelPath)
		}
		ingestedCount++
	}

	if err := tx.Commit(); err != nil {
		return 0, apperror.WrapSimple(err, "failed to commit transaction")
	}

	return ingestedCount, nil
}

// QueryTreeDb executes a search against the Split SQLite database.
func QueryTreeDb(dbPath string, mode SearchMode, pattern string) ([]TreeItem, error) {
	resolved := ResolveTreeDBPath(dbPath)
	if _, err := os.Stat(resolved); os.IsNotExist(err) {
		return nil, apperror.NewNotFoundError("tree database not found at " + resolved + " (run 'gitmap tree-learn -f <file>' first)")
	}

	db, err := OpenTreeDB(resolved)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT TreeFileId, RootPath, RelPath, FileName, DirPath, Extension, SizeBytes, IsDir, Depth, UpdatedAt
		FROM TreeFile
		ORDER BY DirPath, FileName;
	`)
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query TreeFile")
	}
	defer rows.Close()

	all := make([]TreeItem, 0, 256)
	for rows.Next() {
		var item TreeItem
		var isDirInt int
		if err := rows.Scan(
			&item.Id,
			&item.RootPath,
			&item.RelPath,
			&item.FileName,
			&item.DirPath,
			&item.Extension,
			&item.SizeBytes,
			&isDirInt,
			&item.Depth,
			&item.UpdatedAt,
		); err != nil {
			return nil, apperror.WrapSimple(err, "failed to scan row")
		}
		item.IsDir = (isDirInt == 1)
		all = append(all, item)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "error reading rows")
	}

	if strings.TrimSpace(pattern) == "" {
		return all, nil
	}

	filter := TreeSearchFilter{
		Pattern:    pattern,
		SearchMode: mode,
	}
	return FilterTreeItems(all, filter)
}
