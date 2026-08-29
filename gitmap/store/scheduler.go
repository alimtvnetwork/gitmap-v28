package store

import (
	"database/sql"
)

const (
	SQLCreateSchedulerTasksTable = `
CREATE TABLE IF NOT EXISTS scheduler_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    interval_val TEXT,
    delay_val TEXT,
    is_scheduled INTEGER DEFAULT 0,
    has_delay INTEGER DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);`
)

// SchedulerTask represents a scheduled task.
type SchedulerTask struct {
	ID          int
	Name        string
	IntervalVal string
	DelayVal    string
	IsScheduled bool
	HasDelay    bool
	CreatedAt   string
}

// InitSchedulerTable creates the scheduler_tasks table.
func (db *DB) InitSchedulerTable() error {
	_, err := db.conn.Exec(SQLCreateSchedulerTasksTable)
	return err
}

// InsertSchedule adds a new schedule to the database.
func (db *DB) InsertSchedule(name, intervalVal, delayVal string, isScheduled, hasDelay bool) error {
	q := `INSERT INTO scheduler_tasks (name, interval_val, delay_val, is_scheduled, has_delay) VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(q, name, intervalVal, delayVal, isScheduled, hasDelay)
	return err
}

// ListSchedules returns all tasks.
func (db *DB) ListSchedules() ([]SchedulerTask, error) {
	q := `SELECT id, name, interval_val, delay_val, is_scheduled, has_delay, created_at FROM scheduler_tasks`
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
		if err := rows.Scan(&t.ID, &t.Name, &t.IntervalVal, &t.DelayVal, &t.IsScheduled, &t.HasDelay, &t.CreatedAt); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}
