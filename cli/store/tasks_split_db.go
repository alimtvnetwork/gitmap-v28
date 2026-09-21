// Package store — tasks_split_db.go manages the split tasks database.
package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
		_, execErr := conn.Exec(stmt)
		if execErr != nil {
			return apperror.WrapSimple(execErr, "init tasks root schema")
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
