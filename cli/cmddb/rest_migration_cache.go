package cmddb

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// MigrationCache stores DDL execution history and clone records.
type MigrationCache struct {
	mu     sync.Mutex
	dbPath string
	conn   *sql.DB
}

// OpenMigrationCache opens or creates the migration cache SQLite database.
func OpenMigrationCache(customPath ...string) (*MigrationCache, error) {
	path := resolveCachePath(customPath...)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	openRes := store.OpenSQLiteDB(path)
	if openRes.IsFailure() {
		return nil, apperror.WrapSimple(openRes.AppError(), "open migration cache "+path)
	}
	conn := openRes.Value
	conn.SetMaxOpenConns(1)

	if err := initCacheTables(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return &MigrationCache{dbPath: path, conn: conn}, nil
}

func resolveCachePath(custom ...string) string {
	if len(custom) > 0 && custom[0] != "" {
		return custom[0]
	}
	candidates := []string{
		filepath.Join(".gitmap", "repodb", "dbrest_cache.db"),
		filepath.Join(".ai-memory", "dbrest_cache.db"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Dir(c)); err == nil {
			return c
		}
	}

	return candidates[0]
}

func initCacheTables(conn *sql.DB) error {
	ddls := []string{
		`CREATE TABLE IF NOT EXISTS MigrationLog (
			MigrationLogId INTEGER PRIMARY KEY AUTOINCREMENT,
			DatabaseAlias TEXT NOT NULL,
			Version TEXT NOT NULL,
			MigrationName TEXT NOT NULL,
			SqlContent TEXT NOT NULL,
			ExecutedAt TEXT NOT NULL,
			IsSuccess INTEGER NOT NULL DEFAULT 1,
			ErrorMessage TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS AlterationHistory (
			AlterationHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
			DatabaseAlias TEXT NOT NULL,
			TableName TEXT NOT NULL,
			AlterationType TEXT NOT NULL,
			SqlStatement TEXT NOT NULL,
			ExecutedAt TEXT NOT NULL,
			IsSuccess INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS ViewDefinition (
			ViewDefinitionId INTEGER PRIMARY KEY AUTOINCREMENT,
			DatabaseAlias TEXT NOT NULL,
			ViewName TEXT NOT NULL,
			QuerySql TEXT NOT NULL,
			CreatedAt TEXT NOT NULL,
			IsActive INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS ClonedDatabase (
			ClonedDatabaseId INTEGER PRIMARY KEY AUTOINCREMENT,
			SourceAlias TEXT NOT NULL,
			TargetAlias TEXT NOT NULL,
			SourcePath TEXT NOT NULL,
			TargetPath TEXT NOT NULL,
			FileSizeBytes INTEGER NOT NULL DEFAULT 0,
			ClonedAt TEXT NOT NULL,
			IsActive INTEGER NOT NULL DEFAULT 1
		);`,
	}
	for _, ddl := range ddls {
		if _, err := conn.Exec(ddl); err != nil {
			return apperror.WrapSimple(err, "init cache table")
		}
	}

	return nil
}

// LogMigration records a migration execution in MigrationLog.
func (c *MigrationCache) LogMigration(alias, version, name, sqlStr string, isSuccess bool, errMsg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	succInt := 0
	if isSuccess {
		succInt = 1
	}
	q := `INSERT INTO MigrationLog (DatabaseAlias, Version, MigrationName, SqlContent, ExecutedAt, IsSuccess, ErrorMessage) VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := c.conn.Exec(q, alias, version, name, sqlStr, now, succInt, errMsg)

	return err
}

// LogAlteration records a table alteration in AlterationHistory.
func (c *MigrationCache) LogAlteration(alias, tableName, altType, sqlStr string, isSuccess bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	succInt := 0
	if isSuccess {
		succInt = 1
	}
	q := `INSERT INTO AlterationHistory (DatabaseAlias, TableName, AlterationType, SqlStatement, ExecutedAt, IsSuccess) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := c.conn.Exec(q, alias, tableName, altType, sqlStr, now, succInt)

	return err
}

// LogViewDefinition records a view creation in ViewDefinition.
func (c *MigrationCache) LogViewDefinition(alias, viewName, querySql string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	q := `INSERT INTO ViewDefinition (DatabaseAlias, ViewName, QuerySql, CreatedAt, IsActive) VALUES (?, ?, ?, ?, 1)`
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := c.conn.Exec(q, alias, viewName, querySql, now)

	return err
}

// LogClonedDatabase records database clone metadata in ClonedDatabase.
func (c *MigrationCache) LogClonedDatabase(srcAlias, targetAlias, srcPath, targetPath string, size int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	q := `INSERT INTO ClonedDatabase (SourceAlias, TargetAlias, SourcePath, TargetPath, FileSizeBytes, ClonedAt, IsActive) VALUES (?, ?, ?, ?, ?, ?, 1)`
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := c.conn.Exec(q, srcAlias, targetAlias, srcPath, targetPath, size, now)

	return err
}

// Close closes the underlying SQLite cache connection.
func (c *MigrationCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
