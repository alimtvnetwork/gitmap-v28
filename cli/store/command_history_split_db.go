// Package store — command_history_split_db.go manages SQLite Split-DB for command execution history.
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
	sqlCreateCommandHistory = `CREATE TABLE IF NOT EXISTS CommandHistory (
    CommandId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommandLine TEXT NOT NULL,
    CommandName TEXT NOT NULL,
    ExitCode INTEGER NOT NULL DEFAULT 0,
    ExecutedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DurationMs INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_cmd_history_time ON CommandHistory(ExecutedAt DESC);
CREATE INDEX IF NOT EXISTS idx_cmd_history_line ON CommandHistory(CommandLine);`
)

// CommandHistoryEntry models an executed command row.
type CommandHistoryEntry struct {
	ID         int64  `json:"id"`
	Line       string `json:"line"`
	Name       string `json:"name"`
	ExitCode   int    `json:"exitCode"`
	ExecutedAt string `json:"executedAt"`
	DurationMs int64  `json:"durationMs"`
}

// CommandHistorySplitDB manages SQLite storage for command execution history.
type CommandHistorySplitDB struct {
	conn *sql.DB
	path string
}

func resolveCommandHistoryPath(customPath string) string {
	if len(strings.TrimSpace(customPath)) > 0 {
		return customPath
	}
	return filepath.Join(BinaryDataDir(), "history", "commands.db")
}

// OpenCommandHistorySplitDB opens or creates the split-db at customPath or default location.
func OpenCommandHistorySplitDB(customPath string) (*CommandHistorySplitDB, error) {
	dbPath := resolveCommandHistoryPath(customPath)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "mkdir history directory")
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open history sqlite")
	}
	if _, errPragma := conn.Exec("PRAGMA busy_timeout = 3000;"); errPragma != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(errPragma, "set history busy_timeout pragma")
	}
	s := &CommandHistorySplitDB{conn: conn, path: dbPath}
	if err := s.initSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the underlying database connection.
func (s *CommandHistorySplitDB) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *CommandHistorySplitDB) initSchema() error {
	_, err := s.conn.Exec(sqlCreateCommandHistory)
	if err != nil {
		return apperror.WrapSimple(err, "init CommandHistory schema")
	}
	return nil
}

// InsertCommandRecord inserts a new command execution entry.
func (s *CommandHistorySplitDB) InsertCommandRecord(cmdLine, cmdName string, exitCode int, durationMs int64) error {
	if s == nil || s.conn == nil {
		return nil
	}
	cleanLine := strings.TrimSpace(cmdLine)
	if cleanLine == "" {
		return nil
	}
	const q = `INSERT INTO CommandHistory (CommandLine, CommandName, ExitCode, DurationMs) VALUES (?, ?, ?, ?);`
	_, err := s.conn.Exec(q, cleanLine, cmdName, exitCode, durationMs)
	if isFatalDatabaseError(err) {
		return apperror.WrapSimple(err, "insert CommandHistory entry")
	}
	return nil
}

func isFatalDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	return !isDatabaseBusy(err)
}

func isDatabaseBusy(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "busy") || strings.Contains(msg, "locked")
}

// RecordCommandSafely opens the default command history db, inserts the record, and closes it, ignoring errors.
func RecordCommandSafely(cmdLine, cmdName string, exitCode int, durationMs int64) {
	histDB, err := OpenCommandHistorySplitDB("")
	if err != nil {
		return
	}
	defer histDB.Close()
	_ = histDB.InsertCommandRecord(cmdLine, cmdName, exitCode, durationMs)
}

// ListRecentCommands returns up to limit recent command history rows.
func (s *CommandHistorySplitDB) ListRecentCommands(limit int) ([]CommandHistoryEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	const q = `SELECT CommandId, CommandLine, CommandName, ExitCode, ExecutedAt, DurationMs
FROM CommandHistory ORDER BY ExecutedAt DESC, CommandId DESC LIMIT ?;`
	rows, err := s.conn.Query(q, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query CommandHistory")
	}
	defer rows.Close()
	return scanCommandHistoryEntries(rows)
}

func scanCommandHistoryEntries(rows *sql.Rows) ([]CommandHistoryEntry, error) {
	var results []CommandHistoryEntry
	for rows.Next() {
		var e CommandHistoryEntry
		if err := rows.Scan(&e.ID, &e.Line, &e.Name, &e.ExitCode, &e.ExecutedAt, &e.DurationMs); err != nil {
			return nil, apperror.WrapSimple(err, "scan CommandHistory entry")
		}
		results = append(results, e)
	}
	return results, nil
}

// SuggestCommands returns distinct past command lines matching prefix.
func (s *CommandHistorySplitDB) SuggestCommands(prefix string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	const q = `SELECT DISTINCT CommandLine FROM CommandHistory
WHERE CommandLine LIKE ? ORDER BY ExecutedAt DESC LIMIT ?;`
	rows, err := s.conn.Query(q, prefix+"%", limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "suggest CommandHistory")
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var l string
		if err := rows.Scan(&l); err == nil && l != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}

// ClearCommands purges all historical command entries.
func (s *CommandHistorySplitDB) ClearCommands() error {
	const q = `DELETE FROM CommandHistory;`
	_, err := s.conn.Exec(q)
	if err != nil {
		return apperror.WrapSimple(err, "clear CommandHistory")
	}
	return nil
}
