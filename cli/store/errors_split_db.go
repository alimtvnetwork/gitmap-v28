// Package store — errors_split_db.go manages SQLite Split-DB storage for GitMap internal errors.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
	RepoPath      string `json:"repoPath,omitempty"`
}

// ErrorsSplitDB manages the SQLite database connection for internal errors.
type ErrorsSplitDB struct {
	conn *sql.DB
	Path string
}

// ErrorsDBPath returns the canonical path to gitmap-errors.db in the data folder.
func ErrorsDBPath() string {
	dir := BinaryDataDir()
	dbPath := filepath.Join(dir, ErrorsDBFileName)
	if fallback := findGlobalFallbackDB(ErrorsDBFileName, dbPath); fallback != "" {
		return fallback
	}

	if isGlobalUserDataFallbackNeeded(dir) {
		return resolveGlobalUserDataDB(ErrorsDBFileName, dbPath)
	}

	_ = os.MkdirAll(dir, 0755)

	return dbPath
}

func isGlobalUserDataFallbackNeeded(dir string) bool {
	if binaryDataDirOverride != "" {
		return false
	}
	lower := strings.ToLower(dir)
	return strings.Contains(lower, "go-build") || strings.Contains(lower, "\\temp\\") || strings.Contains(lower, "/temp/")
}

func resolveGlobalUserDataDB(fileName, defaultPath string) string {
	global := GlobalUserDataDir()
	if global == "" {
		return defaultPath
	}
	_ = os.MkdirAll(global, 0755)
	return filepath.Join(global, fileName)
}

// OpenErrorsSplitDB opens the canonical split database for internal errors.
func OpenErrorsSplitDB() ErrorsSplitDBWrap {
	return OpenErrorsSplitDBAt(ErrorsDBPath())
}

// OpenErrorsSplitDBAt opens or creates an errors split database at a specific path.
func OpenErrorsSplitDBAt(dbPath string) ErrorsSplitDBWrap {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return result.Fail[*ErrorsSplitDB](apperror.WrapSimple(err, "errors_split.mkdir"))
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return result.Fail[*ErrorsSplitDB](apperror.WrapSimple(err, "errors_split.open"))
	}

	return initErrorsSplitConn(conn, dbPath)
}

func initErrorsSplitConn(conn *sql.DB, dbPath string) ErrorsSplitDBWrap {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return result.Fail[*ErrorsSplitDB](err)
	}

	db := &ErrorsSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return result.Fail[*ErrorsSplitDB](err)
	}

	return result.Ok(db)
}

// InitSchema creates the InternalErrorLog and FailedCommand tables and indexes if absent.
func (s *ErrorsSplitDB) InitSchema() *apperror.AppError {
	if _, err := s.conn.Exec(sqlCreateInternalErrorLog); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_schema")
	}
	if _, err := s.conn.Exec(sqlCreateFailedCommand); err != nil {
		return apperror.WrapSimple(err, "errors_split.init_failed_commands_schema")
	}

	return nil
}

// Close closes the underlying database connection.
func (s *ErrorsSplitDB) Close() *apperror.AppError {
	if s.conn == nil {
		return nil
	}

	if err := s.conn.Close(); err != nil {
		return apperror.WrapSimple(err, "ErrorsSplitDB.Close")
	}

	return nil
}

func (s *ErrorsSplitDB) FederatedClearErrors() *apperror.AppError {
	rows, err := s.conn.Query(`SELECT DISTINCT RepoPath FROM RootErrorIndex WHERE RepoPath IS NOT NULL AND RepoPath != ''`)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.clear.query")
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var p string
		scanErr := rows.Scan(&p)
		if scanErr != nil {
			continue
		}
		paths = append(paths, p)
	}

	for _, p := range paths {
		clearRepoErrorsFile(p)
	}

	_, err = s.conn.Exec(`DELETE FROM RootErrorIndex`)
	if err != nil {
		return apperror.WrapSimple(err, "errors_split.clear")
	}
	return nil
}

func clearRepoErrorsFile(repoPath string) {
	repoDBPath := filepath.Join(repoPath, ".gitmap", "errors.db")
	repoConn, err := sql.Open("sqlite", repoDBPath)
	if err != nil {
		return
	}
	defer repoConn.Close()

	_, errDel := repoConn.Exec(`DELETE FROM RepoErrorDB`)
	if errDel != nil {
		return
	}
}

func (s *ErrorsSplitDB) FederatedGetError(id int64) InternalErrorRecordWrap {
	row := s.conn.QueryRow(`SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex WHERE InternalErrorLogId = ?`, id)
	var rec InternalErrorRecord
	var isResolvedInt int

	err := row.Scan(
		&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message,
		&rec.GitMapVersion, &isResolvedInt, &rec.CreatedAt, &rec.RepoPath,
	)
	if err != nil {
		return result.Fail[*InternalErrorRecord](apperror.WrapSimple(err, "errors_split.get"))
	}

	rec.IsResolved = isResolvedInt > 0
	populateRepoErrorDetails(&rec)
	return result.Ok(&rec)
}

func (s *ErrorsSplitDB) FederatedListErrors(limit int, unresolvedOnly bool) InternalErrorRecordSliceWrap {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex ORDER BY InternalErrorLogId DESC LIMIT ?`
	if unresolvedOnly {
		query = `SELECT InternalErrorLogId, ErrorCode, ErrorType, Command, Message, GitMapVersion, IsResolved, CreatedAt, COALESCE(RepoPath, '') FROM RootErrorIndex WHERE IsResolved = 1 ORDER BY InternalErrorLogId DESC LIMIT ?`
	}

	rows, err := s.conn.Query(query, limit)
	if err != nil {
		return result.FailSlice[InternalErrorRecord](apperror.WrapSimple(err, "errors_split.list"))
	}
	defer rows.Close()

	var records []InternalErrorRecord
	for rows.Next() {
		var rec InternalErrorRecord
		var isResolvedInt int
		scanErr := rows.Scan(&rec.ID, &rec.ErrorCode, &rec.ErrorType, &rec.Command, &rec.Message, &rec.GitMapVersion, &isResolvedInt, &rec.CreatedAt, &rec.RepoPath)
		if scanErr != nil {
			continue
		}
		rec.IsResolved = isResolvedInt > 0
		records = append(records, rec)
	}

	for i := range records {
		populateRepoErrorDetails(&records[i])
	}

	return result.OkSlice(records)
}
