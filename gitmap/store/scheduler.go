package store

import (
	"database/sql"
)

const (
	SQLCreateSchedulerTasksTable = `
CREATE TABLE IF NOT EXISTS scheduler_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    macro_name TEXT,
    command_line TEXT,
    interval_val TEXT,
    delay_val TEXT,
    is_scheduled INTEGER DEFAULT 0,
    has_delay INTEGER DEFAULT 0,
    is_startup INTEGER DEFAULT 0,
    run_count INTEGER DEFAULT 0,
    last_run_at TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);`
)

// SchedulerTask represents a scheduled task.
type SchedulerTask struct {
	ID          int    `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	MacroName   string `json:"macroName,omitempty" yaml:"macroName,omitempty"`
	CommandLine string `json:"commandLine,omitempty" yaml:"commandLine,omitempty"`
	IntervalVal string `json:"interval" yaml:"interval"`
	DelayVal    string `json:"delay,omitempty" yaml:"delay,omitempty"`
	IsScheduled bool   `json:"isScheduled" yaml:"isScheduled"`
	HasDelay    bool   `json:"hasDelay" yaml:"hasDelay"`
	IsStartup   bool   `json:"isStartup" yaml:"isStartup"`
	RunCount    int    `json:"runCount" yaml:"runCount"`
	LastRunAt   string `json:"lastRunAt,omitempty" yaml:"lastRunAt,omitempty"`
	CreatedAt   string `json:"createdAt" yaml:"createdAt"`
}

// InitSchedulerTable creates the scheduler_tasks table and ensures schema columns.
func (db *DB) InitSchedulerTable() error {
	if _, err := db.conn.Exec(SQLCreateSchedulerTasksTable); err != nil {
		return err
	}
	db.migrateSchedulerColumns()
	return nil
}

func (db *DB) migrateSchedulerColumns() {
	_ = db.addTableColumnSafe("scheduler_tasks", "macro_name", "TEXT")
	_ = db.addTableColumnSafe("scheduler_tasks", "command_line", "TEXT")
	_ = db.addTableColumnSafe("scheduler_tasks", "is_startup", "INTEGER DEFAULT 0")
	_ = db.addTableColumnSafe("scheduler_tasks", "run_count", "INTEGER DEFAULT 0")
	_ = db.addTableColumnSafe("scheduler_tasks", "last_run_at", "TEXT")
}

func (db *DB) addTableColumnSafe(tableName, colName, colType string) error {
	q := "ALTER TABLE " + tableName + " ADD COLUMN " + colName + " " + colType
	_, err := db.conn.Exec(q)
	return err
}

// InsertSchedule adds a new schedule to the database.
func (db *DB) InsertSchedule(t SchedulerTask) error {
	q := `INSERT INTO scheduler_tasks (name, macro_name, command_line, interval_val, delay_val, is_scheduled, has_delay, is_startup) 
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	      ON CONFLICT(name) DO UPDATE SET 
	          macro_name=excluded.macro_name,
	          command_line=excluded.command_line,
	          interval_val=excluded.interval_val,
	          delay_val=excluded.delay_val,
	          is_scheduled=excluded.is_scheduled,
	          has_delay=excluded.has_delay,
	          is_startup=excluded.is_startup`
	_, err := db.conn.Exec(q, t.Name, t.MacroName, t.CommandLine, t.IntervalVal, t.DelayVal, t.IsScheduled, t.HasDelay, t.IsStartup)
	return err
}

// GetSchedule retrieves a scheduled task by name.
func (db *DB) GetSchedule(name string) (*SchedulerTask, error) {
	q := `SELECT id, name, COALESCE(macro_name,''), COALESCE(command_line,''), interval_val, delay_val, is_scheduled, has_delay, is_startup, run_count, COALESCE(last_run_at,''), created_at FROM scheduler_tasks WHERE name = ?`
	row := db.conn.QueryRow(q, name)
	var t SchedulerTask
	err := row.Scan(&t.ID, &t.Name, &t.MacroName, &t.CommandLine, &t.IntervalVal, &t.DelayVal, &t.IsScheduled, &t.HasDelay, &t.IsStartup, &t.RunCount, &t.LastRunAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteSchedule removes a schedule by name.
func (db *DB) DeleteSchedule(name string) error {
	q := `DELETE FROM scheduler_tasks WHERE name = ?`
	_, err := db.conn.Exec(q, name)
	return err
}

// UpdateScheduleRun increments run count and records last execution timestamp.
func (db *DB) UpdateScheduleRun(name, timestamp string) error {
	q := `UPDATE scheduler_tasks SET run_count = run_count + 1, last_run_at = ? WHERE name = ?`
	_, err := db.conn.Exec(q, timestamp, name)
	return err
}

// ListSchedules returns all tasks.
func (db *DB) ListSchedules() ([]SchedulerTask, error) {
	q := `SELECT id, name, COALESCE(macro_name,''), COALESCE(command_line,''), interval_val, delay_val, is_scheduled, has_delay, is_startup, run_count, COALESCE(last_run_at,''), created_at FROM scheduler_tasks ORDER BY id ASC`
	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseScheduleRows(rows), nil
}

// parseScheduleRows parses rows into tasks.
func parseScheduleRows(rows *sql.Rows) []SchedulerTask {
	var tasks []SchedulerTask
	for rows.Next() {
		var t SchedulerTask
		if err := rows.Scan(&t.ID, &t.Name, &t.MacroName, &t.CommandLine, &t.IntervalVal, &t.DelayVal, &t.IsScheduled, &t.HasDelay, &t.IsStartup, &t.RunCount, &t.LastRunAt, &t.CreatedAt); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}
