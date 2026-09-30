// Package store — tasks_split_db.go manages the split tasks database.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	_ "modernc.org/sqlite"
)

const (
	sqlCreateTaskQueue = `CREATE TABLE IF NOT EXISTS TaskQueue (
    TaskQueueId    INTEGER PRIMARY KEY AUTOINCREMENT,
    QueueId        TEXT NOT NULL UNIQUE,
    Section        TEXT NOT NULL DEFAULT '',
    Action         TEXT NOT NULL,
    Target         TEXT NOT NULL,
    ForwardPayload TEXT NOT NULL DEFAULT '',
    InversePayload TEXT NOT NULL DEFAULT '',
    Status         TEXT NOT NULL DEFAULT 'pending',
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxTaskQueue_Section ON TaskQueue(Section);
CREATE INDEX IF NOT EXISTS IdxTaskQueue_Status ON TaskQueue(Status);`

	sqlCreateTaskHistory = `CREATE TABLE IF NOT EXISTS TaskHistory (
    TaskHistoryId  INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskId         TEXT NOT NULL UNIQUE,
    Section        TEXT NOT NULL DEFAULT '',
    Action         TEXT NOT NULL,
    Target         TEXT NOT NULL,
    ForwardPayload TEXT NOT NULL DEFAULT '',
    InversePayload TEXT NOT NULL DEFAULT '',
    Status         TEXT NOT NULL DEFAULT 'completed',
    RestoredAt     TEXT NOT NULL DEFAULT '',
    ExecutedAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt      TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxTaskHistory_Section ON TaskHistory(Section);
CREATE INDEX IF NOT EXISTS IdxTaskHistory_Action ON TaskHistory(Action);
CREATE INDEX IF NOT EXISTS IdxTaskHistory_CreatedAt ON TaskHistory(CreatedAt);`
)

// TasksSplitDB represents an isolated SQLite database connection for tasks.
type TasksSplitDB struct {
	*DB
	Path string
}

// TasksRootDBPath returns the full path to the tasks root SQLite DB.
func TasksRootDBPath() string {
	return ResolveTasksRootDbPath("")
}

// OpenTasksRootSplitDB opens or initializes the tasks root split database.
func OpenTasksRootSplitDB() (*TasksSplitDB, error) {
	return OpenTasksRootSplitDBAt(TasksRootDBPath())
}

// OpenTasksRootSplitDBAt opens or initializes the tasks root split database at a path.
func OpenTasksRootSplitDBAt(dbPath string) (*TasksSplitDB, error) {
	mkErr := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if mkErr != nil {
		return nil, apperror.WrapSimple(mkErr, "open tasks root split db: mkdir")
	}

	innerDB, openErr := openDBAt(dbPath)
	if openErr != nil {
		return nil, apperror.WrapSimple(openErr, "open tasks root split db: open")
	}

	return initTasksSplitConn(innerDB, dbPath)
}

func initTasksSplitConn(innerDB *DB, dbPath string) (*TasksSplitDB, error) {
	schemaErr := initTasksSchema(innerDB.Conn())
	if schemaErr != nil {
		_ = innerDB.Close()
		return nil, schemaErr
	}

	tasksDB := &TasksSplitDB{DB: innerDB, Path: dbPath}
	tasksDB.registerWithRegistry()

	return tasksDB, nil
}

func initTasksSchema(conn *sql.DB) error {
	statements := []string{
		constants.SQLCreateTaskType,
		constants.SQLCreatePendingTask,
		constants.SQLCreateCompletedTask,
		constants.SQLSeedTaskTypes,
		sqlCreateTaskQueue,
		sqlCreateTaskHistory,
	}
	for _, stmt := range statements {
		res := ExecWrapper(conn, stmt)
		if res.IsFailure {
			return apperror.WrapSimple(res.Error, "init tasks root schema")
		}
	}
	return nil
}

func (db *TasksSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}
	defer master.Close()

	entry := SplitDatabaseEntry{
		DatabaseType:  "tasks",
		DatabaseKey:   "tasks_root",
		DatabasePath:  db.Path,
		Status:        "active",
		IsActive:      true,
		Description:   "Root database for tasks, queue, and execution history",
		SchemaVersion: 1,
	}
	_ = master.RegisterSplitDB(entry)
}

// Close terminates the database connection.
func (db *TasksSplitDB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}

// Conn returns the underlying SQLite connection.
func (db *TasksSplitDB) Conn() *sql.DB {
	return db.DB.Conn()
}

// InsertTaskHistory adds an execution audit record into TaskHistory.
func (db *DB) InsertTaskHistory(
	taskId,
	section,
	action,
	target,
	forward,
	inverse,
	status string,
) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO TaskHistory (TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status, ExecutedAt, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := ExecWrapper(db.Conn(), query, taskId, section, action, target, forward, inverse, status, now, now)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "InsertTaskHistory")
	}

	return nil
}

// ListTaskHistory queries audit records with section filter and pagination.
func (db *DB) ListTaskHistory(section string, limit, offset int) ([]model.TaskHistoryRecord, error) {
	boundedLimit, boundedOffset := sanitizeTaskHistoryBounds(limit, offset)
	rows, err := db.queryTaskHistoryRows(section, boundedLimit, boundedOffset)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListTaskHistory query")
	}
	defer rows.Close()

	return scanTaskHistoryRows(rows)
}

func sanitizeTaskHistoryBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func (db *DB) queryTaskHistoryRows(section string, limit, offset int) (*sql.Rows, error) {
	if section != "" && section != "all" {
		query := `SELECT TaskHistoryId, TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status, RestoredAt, ExecutedAt, CreatedAt
FROM TaskHistory WHERE Section = ? ORDER BY TaskHistoryId DESC LIMIT ? OFFSET ?`
		return db.Conn().Query(query, section, limit, offset)
	}

	query := `SELECT TaskHistoryId, TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status, RestoredAt, ExecutedAt, CreatedAt
FROM TaskHistory ORDER BY TaskHistoryId DESC LIMIT ? OFFSET ?`
	return db.Conn().Query(query, limit, offset)
}

func scanTaskHistoryRows(rows *sql.Rows) ([]model.TaskHistoryRecord, error) {
	var records []model.TaskHistoryRecord
	for rows.Next() {
		var r model.TaskHistoryRecord
		scanErr := rows.Scan(
			&r.TaskHistoryId, &r.TaskId, &r.Section, &r.Action, &r.Target,
			&r.ForwardPayload, &r.InversePayload, &r.Status,
			&r.RestoredAt, &r.ExecutedAt, &r.CreatedAt,
		)
		if scanErr == nil {
			records = append(records, r)
		}
	}

	return records, nil
}
