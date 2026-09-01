// Package store — schedule_split_db.go: per-schedule isolated SQLite split databases
// for storing full schedule configuration and massive execution run logs.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

// ScheduleSplitDB wraps an isolated SQLite database connection for a single schedule.
type ScheduleSplitDB struct {
	conn *sql.DB
	Slug string
	Path string
}

// ScheduleRunRecord represents a single recorded execution of a scheduled task.
type ScheduleRunRecord struct {
	ID          int64  `json:"id" yaml:"id"`
	RunNumber   int    `json:"runNumber" yaml:"runNumber"`
	TriggerType string `json:"triggerType" yaml:"triggerType"`
	RunnerUser  string `json:"runnerUser" yaml:"runnerUser"`
	StartedAt   string `json:"startedAt" yaml:"startedAt"`
	FinishedAt  string `json:"finishedAt" yaml:"finishedAt"`
	DurationMS  int64  `json:"durationMs" yaml:"durationMs"`
	IsSuccess   bool   `json:"isSuccess" yaml:"isSuccess"`
	ExitCode    int    `json:"exitCode" yaml:"exitCode"`
	Output      string `json:"output" yaml:"output"`
	ErrorMsg    string `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`
	CreatedAt   string `json:"createdAt" yaml:"createdAt"`
}

// ScheduleConfig represents metadata stored in the split schedule database.
type ScheduleConfig struct {
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug" yaml:"slug"`
	MacroName   string `json:"macroName,omitempty" yaml:"macroName,omitempty"`
	CommandLine string `json:"commandLine,omitempty" yaml:"commandLine,omitempty"`
	IntervalVal string `json:"interval" yaml:"interval"`
	DelayVal    string `json:"delay,omitempty" yaml:"delay,omitempty"`
	IsEnabled   bool   `json:"isEnabled" yaml:"isEnabled"`
	IsStartup   bool   `json:"isStartup" yaml:"isStartup"`
	CreatedAt   string `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   string `json:"updatedAt" yaml:"updatedAt"`
}

// ScheduleSlug converts a schedule name to a clean, lowercase filesystem slug.
func ScheduleSlug(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	reg := regexp.MustCompile(`[^a-z0-9_-]+`)
	slug := reg.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "schedule-default"
	}
	return slug
}

// ScheduleDBDir returns the dedicated schedules subfolder within BinaryDataDir.
func ScheduleDBDir() string {
	dir := filepath.Join(BinaryDataDir(), "schedules")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// ScheduleDBPath returns the full path to a schedule's isolated split SQLite DB.
func ScheduleDBPath(slug string) string {
	return filepath.Join(ScheduleDBDir(), slug+".db")
}

// OpenScheduleSplitDB opens or creates the isolated SQLite DB for a schedule.
func OpenScheduleSplitDB(slug string) (*ScheduleSplitDB, error) {
	dbPath := ScheduleDBPath(slug)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open schedule split db "+slug)
	}
	s := &ScheduleSplitDB{conn: conn, Slug: slug, Path: dbPath}
	if err := s.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return s, nil
}

const (
	sqlCreateScheduleConfig = `
CREATE TABLE IF NOT EXISTS schedule_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    slug TEXT NOT NULL,
    macro_name TEXT,
    command_line TEXT,
    interval_val TEXT,
    delay_val TEXT,
    is_enabled INTEGER DEFAULT 1,
    is_startup INTEGER DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

	sqlCreateScheduleLogs = `
CREATE TABLE IF NOT EXISTS schedule_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_number INTEGER NOT NULL,
    trigger_type TEXT DEFAULT 'scheduled',
    runner_user TEXT,
    started_at TEXT NOT NULL,
    finished_at TEXT NOT NULL,
    duration_ms INTEGER DEFAULT 0,
    is_success INTEGER DEFAULT 1,
    exit_code INTEGER DEFAULT 0,
    output TEXT,
    error_msg TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);`
)

// InitSchema creates tables and indexes in the split DB.
func (s *ScheduleSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreateScheduleConfig); err != nil {
		return apperror.WrapSimple(err, "init schedule_config table")
	}
	if _, err := s.conn.Exec(sqlCreateScheduleLogs); err != nil {
		return apperror.WrapSimple(err, "init schedule_logs table")
	}
	return nil
}

// SaveConfig persists or updates the schedule metadata in the split DB.
func (s *ScheduleSplitDB) SaveConfig(cfg ScheduleConfig) error {
	q := `INSERT INTO schedule_config (name, slug, macro_name, command_line, interval_val, delay_val, is_enabled, is_startup, updated_at)
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	      ON CONFLICT(name) DO UPDATE SET
	          macro_name=excluded.macro_name,
	          command_line=excluded.command_line,
	          interval_val=excluded.interval_val,
	          delay_val=excluded.delay_val,
	          is_enabled=excluded.is_enabled,
	          is_startup=excluded.is_startup,
	          updated_at=CURRENT_TIMESTAMP`
	_, err := s.conn.Exec(q, cfg.Name, cfg.Slug, cfg.MacroName, cfg.CommandLine, cfg.IntervalVal, cfg.DelayVal, cfg.IsEnabled, cfg.IsStartup)
	return err
}

// RecordRun records a single execution event in the schedule logs table.
func (s *ScheduleSplitDB) RecordRun(r ScheduleRunRecord) error {
	q := `INSERT INTO schedule_logs (run_number, trigger_type, runner_user, started_at, finished_at, duration_ms, is_success, exit_code, output, error_msg)
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.conn.Exec(q, r.RunNumber, r.TriggerType, r.RunnerUser, r.StartedAt, r.FinishedAt, r.DurationMS, r.IsSuccess, r.ExitCode, r.Output, r.ErrorMsg)
	return err
}

// GetRuns fetches execution records from the split DB in descending order.
func (s *ScheduleSplitDB) GetRuns(limit int) ([]ScheduleRunRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	q := `SELECT id, run_number, trigger_type, runner_user, started_at, finished_at, duration_ms, is_success, exit_code, output, error_msg, created_at
	      FROM schedule_logs ORDER BY id DESC LIMIT ?`
	rows, err := s.conn.Query(q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseRunRows(rows), nil
}

func parseRunRows(rows *sql.Rows) []ScheduleRunRecord {
	var list []ScheduleRunRecord
	for rows.Next() {
		var r ScheduleRunRecord
		var isSuccessInt int
		err := rows.Scan(&r.ID, &r.RunNumber, &r.TriggerType, &r.RunnerUser, &r.StartedAt, &r.FinishedAt, &r.DurationMS, &isSuccessInt, &r.ExitCode, &r.Output, &r.ErrorMsg, &r.CreatedAt)
		if err == nil {
			r.IsSuccess = isSuccessInt == 1
			list = append(list, r)
		}
	}
	return list
}

// ResetLogs truncates the schedule logs table in the split DB.
func (s *ScheduleSplitDB) ResetLogs() error {
	_, err := s.conn.Exec("DELETE FROM schedule_logs")
	return err
}

// Close closes the SQLite connection.
func (s *ScheduleSplitDB) Close() error {
	return s.conn.Close()
}

// DeleteScheduleSplitDB removes the SQLite database file for a schedule.
func DeleteScheduleSplitDB(slug string) error {
	path := ScheduleDBPath(slug)
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

// MigrateAllScheduleSplitDBs iterates through all split databases and migrates schema.
func MigrateAllScheduleSplitDBs() error {
	dir := ScheduleDBDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".db") {
			migrateSingleSplitDB(strings.TrimSuffix(e.Name(), ".db"))
		}
	}
	return nil
}

func migrateSingleSplitDB(slug string) {
	db, err := OpenScheduleSplitDB(slug)
	if err == nil && db != nil {
		_ = db.InitSchema()
		_ = db.Close()
	}
}
