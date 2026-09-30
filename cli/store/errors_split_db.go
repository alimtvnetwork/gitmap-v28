// Package store — errors_split_db.go manages SQLite Split-DB storage for GitMap internal errors.
package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	ErrorsDBFileName = "gitmap-errors.db"

	sqlCreateInternalErrorLog = `CREATE TABLE IF NOT EXISTS InternalErrorLog (
    InternalErrorLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    ErrorCode          TEXT NOT NULL DEFAULT '',
    ErrorType          TEXT NOT NULL DEFAULT 'general',
    Command            TEXT NOT NULL DEFAULT '',
    Message            TEXT NOT NULL,
    Details            TEXT NULL,
    SourceFile         TEXT NULL,
    ContextJson        TEXT NULL,
    StackTrace         TEXT NULL,
    GitMapVersion      TEXT NOT NULL DEFAULT '',
    IsResolved         INTEGER NOT NULL DEFAULT 0,
    Notes              TEXT NULL,
    Comments           TEXT NULL,
    CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_CreatedAt ON InternalErrorLog(CreatedAt DESC);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_ErrorType ON InternalErrorLog(ErrorType);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_ErrorCode ON InternalErrorLog(ErrorCode);
CREATE INDEX IF NOT EXISTS IdxInternalErrorLog_IsResolved ON InternalErrorLog(IsResolved);`

	sqlCreateFailedCommand = `CREATE TABLE IF NOT EXISTS FailedCommand (
    FailedCommandId INTEGER PRIMARY KEY AUTOINCREMENT,
    Command         TEXT NOT NULL DEFAULT '',
    FullArgs        TEXT NOT NULL DEFAULT '',
    Domain          TEXT NOT NULL DEFAULT 'root',
    ErrorCode       TEXT NOT NULL DEFAULT 'E1001',
    Message         TEXT NOT NULL DEFAULT '',
    Suggestions     TEXT NOT NULL DEFAULT '',
    HitCount        INTEGER NOT NULL DEFAULT 1,
    WorkingDir      TEXT NOT NULL DEFAULT '',
    GitMapVersion   TEXT NOT NULL DEFAULT '',
    IsResolved      INTEGER NOT NULL DEFAULT 0,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LastSeenAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxFailedCommand_Command ON FailedCommand(Command, Domain);
CREATE INDEX IF NOT EXISTS IdxFailedCommand_HitCount ON FailedCommand(HitCount DESC);
CREATE INDEX IF NOT EXISTS IdxFailedCommand_LastSeenAt ON FailedCommand(LastSeenAt DESC);
CREATE VIEW IF NOT EXISTS failed_commands AS SELECT * FROM FailedCommand;
CREATE VIEW IF NOT EXISTS failed_to_detect_commands AS SELECT * FROM FailedCommand;
CREATE VIEW IF NOT EXISTS FailedToDetectCommand AS SELECT * FROM FailedCommand;`
)

// InternalErrorRecord represents a recorded internal error log entry.
type InternalErrorRecord struct {
	ID            int64  `json:"id"`
	ErrorCode     string `json:"errorCode"`
	ErrorType     string `json:"errorType"`
	Command       string `json:"command"`
	Message       string `json:"message"`
	Details       string `json:"details,omitempty"`
	SourceFile    string `json:"sourceFile,omitempty"`
	ContextJson   string `json:"contextJson,omitempty"`
	StackTrace    string `json:"stackTrace,omitempty"`
	GitMapVersion string `json:"gitMapVersion,omitempty"`
	IsResolved    bool   `json:"isResolved"`
	Notes         string `json:"notes,omitempty"`
	Comments      string `json:"comments,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

// ErrorsSplitDB manages the SQLite database connection for internal errors.
type ErrorsSplitDB struct {
	conn *sql.DB
	Path string
}

// ErrorsDBPath returns the canonical path to gitmap-errors.db in the data folder.
func ErrorsDBPath() string {
	dir := BinaryDataDir()
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, ErrorsDBFileName)
}

// OpenErrorsSplitDB opens the canonical split database for internal errors.
func OpenErrorsSplitDB() (*ErrorsSplitDB, error) {
	return OpenErrorsSplitDBAt(ErrorsDBPath())
}

// OpenErrorsSplitDBAt opens or creates an errors split database at a specific path.
func OpenErrorsSplitDBAt(dbPath string) (*ErrorsSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "errors_split.open")
	}

	return initErrorsSplitConn(conn, dbPath)
}

func initErrorsSplitConn(conn *sql.DB, dbPath string) (*ErrorsSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "errors_split.config")
	}

	db := &ErrorsSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return db, nil
}

// InitSchema creates the InternalErrorLog and FailedCommand tables and indexes if absent.
func (s *ErrorsSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreateInternalErrorLog); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_schema")
	}
	if _, err := s.conn.Exec(sqlCreateFailedCommand); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_failed_commands_schema")
	}

	return nil
}

// Close closes the underlying database connection.
func (s *ErrorsSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}
