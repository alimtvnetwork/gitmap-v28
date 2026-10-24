package cmddb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type sysAliasDef struct {
	name    string
	getPath func() string
	dbType  string
	desc    string
}

// InitAliasRegistry ensures the DatabaseAlias table exists in gitmap.db and seeds system aliases.
func InitAliasRegistry() error {
	dbPath := store.DefaultDBPath()
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return openRes.AppError()
	}
	conn := openRes.Value
	defer conn.Close()

	ddl := `CREATE TABLE IF NOT EXISTS DatabaseAlias (
		DatabaseAliasId INTEGER PRIMARY KEY AUTOINCREMENT,
		AliasName TEXT NOT NULL UNIQUE,
		DatabasePath TEXT NOT NULL,
		DatabaseType TEXT NOT NULL DEFAULT 'custom',
		IsActive INTEGER NOT NULL DEFAULT 1,
		FileSizeBytes INTEGER NOT NULL DEFAULT 0,
		SearchHistoryCount INTEGER NOT NULL DEFAULT 0,
		LastAccessedAt TEXT,
		Description TEXT,
		Notes TEXT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_database_alias_name ON DatabaseAlias (AliasName);`

	_, err := conn.Exec(ddl)
	if err != nil {
		return err
	}

	seedSystemAliases(conn)
	return nil
}

func seedSystemAliases(conn *sql.DB) {
	for _, def := range getSystemAliasDefs() {
		seedOneSystemAlias(conn, def)
	}
}

func seedOneSystemAlias(conn *sql.DB, def sysAliasDef) {
	path := def.getPath()
	size := getFileSize(path)
	query := `INSERT INTO DatabaseAlias (AliasName, DatabasePath, DatabaseType, IsActive, FileSizeBytes, SearchHistoryCount, LastAccessedAt, Description, Notes)
VALUES (?, ?, ?, 1, ?, 0, '', ?, '')
ON CONFLICT(AliasName) DO UPDATE SET DatabasePath=excluded.DatabasePath, DatabaseType=excluded.DatabaseType, Description=excluded.Description, FileSizeBytes=CASE WHEN excluded.FileSizeBytes > 0 THEN excluded.FileSizeBytes ELSE DatabaseAlias.FileSizeBytes END`

	if _, err := conn.Exec(query, def.name, path, def.dbType, size, def.desc); err != nil {
		return
	}
}

func getSystemAliasDefs() []sysAliasDef {
	return []sysAliasDef{
		{"master", store.DefaultDBPath, "system-master", "Primary master repository and configuration database"},
		{"pipeline", func() string { return pipelinedb.ResolvePipelineDbPath("") }, "system-pipeline", "CI/CD pipeline runs and telemetry database"},
		{"heatmap", func() string { return filepath.Join(store.EffectiveDataDir(), "heatmap.db") }, "system-heatmap", "AI agent heatmap and coordination database"},
		{"tasks", func() string { return store.ResolveMasterAgentDbPath("") }, "system-tasks", "Master agent swarm task queue database"},
	}
}

// RegisterAlias registers or updates a database alias mapping with file size stat.
func RegisterAlias(name, path, desc string) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return apperror.NewValidationError("alias name cannot be empty")
	}
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return apperror.NewValidationError("database path cannot be empty")
	}

	if err := InitAliasRegistry(); err != nil {
		return err
	}

	cleanPath := filepath.Clean(trimmedPath)
	size := getFileSize(cleanPath)
	dbType := inferDatabaseType(trimmedName, cleanPath)

	return executeRegisterAlias(trimmedName, cleanPath, dbType, desc, size)
}

func executeRegisterAlias(name, path, dbType, desc string, size int64) error {
	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return openRes.AppError()
	}
	conn := openRes.Value
	defer conn.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO DatabaseAlias (AliasName, DatabasePath, DatabaseType, IsActive, FileSizeBytes, SearchHistoryCount, LastAccessedAt, Description, Notes)
VALUES (?, ?, ?, 1, ?, 0, ?, ?, '')
ON CONFLICT(AliasName) DO UPDATE SET DatabasePath=excluded.DatabasePath, DatabaseType=excluded.DatabaseType, FileSizeBytes=excluded.FileSizeBytes, Description=CASE WHEN excluded.Description != '' THEN excluded.Description ELSE DatabaseAlias.Description END, LastAccessedAt=excluded.LastAccessedAt, IsActive=1`

	_, err := conn.Exec(query, name, path, dbType, size, now, desc)
	return err
}

func inferDatabaseType(name, path string) string {
	lowerName := strings.ToLower(name)
	lowerPath := strings.ToLower(path)
	if lowerName == "master" || lowerName == "gitmap" {
		return "system-master"
	}
	if strings.Contains(lowerName, "pipeline") || strings.Contains(lowerPath, "pipeline") {
		return "system-pipeline"
	}
	if strings.Contains(lowerName, "heatmap") || strings.Contains(lowerPath, "heatmap") {
		return "system-heatmap"
	}
	if strings.Contains(lowerName, "task") || strings.Contains(lowerName, "agent") {
		return "system-tasks"
	}
	if strings.Contains(lowerPath, "repo_search") {
		return "split-db"
	}
	return "custom"
}

// ResolveDbPathWithRegistry resolves system aliases, custom aliases, or filesystem paths.
func ResolveDbPathWithRegistry(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if path, ok := resolveSystemDbAlias(trimmed); ok {
		return path, nil
	}
	if path, ok := lookupAliasInRegistry(trimmed); ok {
		return path, nil
	}
	return filepath.Clean(trimmed), nil
}

func resolveSystemDbAlias(input string) (string, bool) {
	if input == "" || strings.EqualFold(input, "master") || strings.EqualFold(input, "gitmap") {
		return store.DefaultDBPath(), true
	}
	if strings.EqualFold(input, "heatmap") {
		return filepath.Join(store.EffectiveDataDir(), "heatmap.db"), true
	}
	if strings.EqualFold(input, "pipeline") || strings.EqualFold(input, "pipelines") {
		return pipelinedb.ResolvePipelineDbPath(""), true
	}
	if strings.EqualFold(input, "tasks") || strings.EqualFold(input, "agents") {
		return store.ResolveMasterAgentDbPath(""), true
	}
	return "", false
}

func lookupAliasInRegistry(name string) (string, bool) {
	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return "", false
	}
	conn := openRes.Value
	defer conn.Close()

	var resolvedPath string
	var isActiveInt int
	query := `SELECT DatabasePath, IsActive FROM DatabaseAlias WHERE LOWER(AliasName) = LOWER(?) LIMIT 1`
	err := conn.QueryRow(query, name).Scan(&resolvedPath, &isActiveInt)
	if err != nil || isActiveInt == 0 {
		return "", false
	}
	return filepath.Clean(resolvedPath), true
}

// IncrementAliasSearchCount increments the search count and updates the last accessed timestamp for an alias or path.
func IncrementAliasSearchCount(aliasOrPath string) {
	trimmed := strings.TrimSpace(aliasOrPath)
	if trimmed == "" {
		return
	}

	cleanPath := filepath.Clean(trimmed)
	now := time.Now().UTC().Format(time.RFC3339)

	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return
	}
	conn := openRes.Value
	defer conn.Close()

	query := `UPDATE DatabaseAlias SET SearchHistoryCount = SearchHistoryCount + 1, LastAccessedAt = ? WHERE LOWER(AliasName) = LOWER(?) OR DatabasePath = ?`
	if _, err := conn.Exec(query, now, trimmed, cleanPath); err != nil {
		return
	}
}

// ListAliases returns all registered database aliases from the DatabaseAlias table.
func ListAliases() ([]DatabaseAlias, error) {
	if err := InitAliasRegistry(); err != nil {
		return nil, err
	}

	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return nil, openRes.AppError()
	}
	conn := openRes.Value
	defer conn.Close()

	query := `SELECT DatabaseAliasId, AliasName, DatabasePath, DatabaseType, IsActive, FileSizeBytes, SearchHistoryCount, COALESCE(LastAccessedAt, ''), COALESCE(Description, ''), COALESCE(Notes, '') FROM DatabaseAlias ORDER BY DatabaseAliasId ASC`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAliasRows(rows)
}

func scanAliasRows(rows *sql.Rows) ([]DatabaseAlias, error) {
	var aliases []DatabaseAlias
	var scanErr error
	for rows.Next() && scanErr == nil {
		scanErr = appendScannedAlias(rows, &aliases)
	}
	if scanErr != nil {
		return nil, scanErr
	}
	return aliases, rows.Err()
}

func appendScannedAlias(rows *sql.Rows, list *[]DatabaseAlias) error {
	var a DatabaseAlias
	var isActiveInt int
	err := rows.Scan(
		&a.DatabaseAliasId, &a.AliasName, &a.DatabasePath, &a.DatabaseType,
		&isActiveInt, &a.FileSizeBytes, &a.SearchHistoryCount,
		&a.LastAccessedAt, &a.Description, &a.Notes,
	)
	if err != nil {
		return err
	}
	a.IsActive = (isActiveInt == 1)
	*list = append(*list, a)
	return nil
}

func getFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fi.Size()
}
