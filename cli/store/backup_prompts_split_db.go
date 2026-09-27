// Package store — backup_prompts_split_db.go manages SQLite Split-DB connections and schemas for running prompts backup.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	sqlCreatePromptBackupBatch = `CREATE TABLE IF NOT EXISTS PromptBackupBatch (
    BatchId TEXT PRIMARY KEY,
    SourcePath TEXT NOT NULL DEFAULT '',
    TotalPrompts INTEGER NOT NULL DEFAULT 0,
    RunningCount INTEGER NOT NULL DEFAULT 0,
    EnqueuedCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Note TEXT NOT NULL DEFAULT ''
);`

	sqlCreatePromptBackupItem = `CREATE TABLE IF NOT EXISTS PromptBackupItem (
    ItemId TEXT PRIMARY KEY,
    BatchId TEXT NOT NULL,
    ProjectName TEXT NOT NULL,
    ProjectPath TEXT NOT NULL,
    ProjectId TEXT NOT NULL,
    ConversationId TEXT NOT NULL,
    SequenceId TEXT NOT NULL DEFAULT '',
    PromptText TEXT NOT NULL,
    PromptStatus TEXT NOT NULL,
    WordCount INTEGER NOT NULL DEFAULT 0,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_batch ON PromptBackupItem(BatchId);
CREATE INDEX IF NOT EXISTS idx_prompt_backup_item_project ON PromptBackupItem(ProjectId);`

	sqlCreatePromptRestoreLedger = `CREATE TABLE IF NOT EXISTS PromptRestoreLedger (
    RestoreId TEXT PRIMARY KEY,
    BatchId TEXT NOT NULL,
    RestoredAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    RestoredCount INTEGER NOT NULL DEFAULT 0,
    IsKept INTEGER NOT NULL DEFAULT 0,
    PrunedAt TEXT NOT NULL DEFAULT '',
    FOREIGN KEY(BatchId) REFERENCES PromptBackupBatch(BatchId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_prompt_restore_batch ON PromptRestoreLedger(BatchId);`
)

// BackupPromptsSplitDB manages dedicated SQLite storage for prompt snapshots.
type BackupPromptsSplitDB struct {
	conn *sql.DB
	path string
}

func resolveBackupPromptsPath(customPath string) string {
	if len(strings.TrimSpace(customPath)) > 0 {
		return customPath
	}
	return filepath.Join(BinaryDataDir(), "backup-prompts", "sql.db")
}

// OpenBackupPromptsSplitDB opens or creates the split-db at customPath or default location.
func OpenBackupPromptsSplitDB(customPath string) (*BackupPromptsSplitDB, error) {
	dbPath := resolveBackupPromptsPath(customPath)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "mkdir backup-prompts")
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open backup-prompts sqlite")
	}
	return initBackupPromptsConn(conn, dbPath)
}

func initBackupPromptsConn(conn *sql.DB, dbPath string) (*BackupPromptsSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "configure backup-prompts sqlite")
	}
	db := &BackupPromptsSplitDB{conn: conn, path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return db, nil
}

// InitSchema creates the backup batch, items, and restore ledger tables.
func (db *BackupPromptsSplitDB) InitSchema() error {
	stmts := []string{
		sqlCreatePromptBackupBatch,
		sqlCreatePromptBackupItem,
		sqlCreatePromptRestoreLedger,
	}
	for _, stmt := range stmts {
		if _, err := db.conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "init backup-prompts schema")
		}
	}
	return nil
}

// Close closes the underlying SQLite database connection.
func (db *BackupPromptsSplitDB) Close() error {
	if db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// Conn returns the raw database connection.
func (db *BackupPromptsSplitDB) Conn() *sql.DB {
	return db.conn
}

// Path returns the database file path.
func (db *BackupPromptsSplitDB) Path() string {
	return db.path
}

// GetStorageInfo returns the database file path and size in bytes.
func (db *BackupPromptsSplitDB) GetStorageInfo() (string, int64, error) {
	info, err := os.Stat(db.path)
	if err == nil {
		return db.path, info.Size(), nil
	}
	if os.IsNotExist(err) {
		return db.path, 0, nil
	}
	return db.path, 0, apperror.WrapSimple(err, "stat backup-prompts db")
}
