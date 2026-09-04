package applogger

import (
	"database/sql"
	"encoding/json"
	"sync"
)

const createLogsTableSQL = `CREATE TABLE IF NOT EXISTS app_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp TEXT NOT NULL,
	level TEXT NOT NULL,
	message TEXT NOT NULL,
	caller TEXT,
	fields_json TEXT,
	stack_trace TEXT
);`

// SQLiteSink writes structured log entries to an SQLite table.
type SQLiteSink struct {
	mu sync.Mutex
	db *sql.DB
}

// NewSQLiteSink creates and initializes the SQLite logging table.
func NewSQLiteSink(db *sql.DB) (*SQLiteSink, error) {
	sink := &SQLiteSink{db: db}
	if err := sink.initTable(); err != nil {
		return nil, err
	}

	return sink, nil
}

// initTable ensures the log table exists.
func (ss *SQLiteSink) initTable() error {
	if ss.db == nil {
		return nil
	}

	_, err := ss.db.Exec(createLogsTableSQL)

	return err
}

// WriteEntry persists a structured log entry to the SQLite database.
func (ss *SQLiteSink) WriteEntry(e LogEntry) error {
	if ss.db == nil {
		return nil
	}

	ss.mu.Lock()
	defer ss.mu.Unlock()

	fieldsJSON, _ := json.Marshal(e.Fields)
	query := `INSERT INTO app_logs (timestamp, level, message, caller, fields_json, stack_trace) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := ss.db.Exec(query, e.Timestamp, e.Level.Name(), e.Message, e.Caller, string(fieldsJSON), e.Stack)

	return err
}

// Sync is a no-op for SQLite transactions.
func (ss *SQLiteSink) Sync() error { return nil }

// Close closes the underlying db connection if needed.
func (ss *SQLiteSink) Close() error { return nil }
