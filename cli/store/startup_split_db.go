// Package store — startup_split_db.go: isolated SQLite database for startup items and logs.
package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	_ "modernc.org/sqlite"
)

// StartupSplitDB wraps the isolated SQLite database connection for startup items.
type StartupSplitDB struct {
	conn *sql.DB
	Path string
}

// StartupItemRecord represents an item configured to run at OS startup or login.
type StartupItemRecord struct {
	StartupItemId int64  `json:"startupItemId"`
	Name          string `json:"name"`
	TargetType    string `json:"targetType"`
	TargetPath    string `json:"targetPath"`
	CommandArgs   string `json:"commandArgs,omitempty"`
	IconPath      string `json:"iconPath,omitempty"`
	RunFrequency  string `json:"runFrequency"`
	IsActive      bool   `json:"isActive"`
	Description   string `json:"description,omitempty"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

// StartupLogRecord represents an execution log for a startup run.
type StartupLogRecord struct {
	StartupLogId  int64  `json:"startupLogId"`
	StartupItemId int64  `json:"startupItemId"`
	RunAt         int64  `json:"runAt"`
	DurationMs    int64  `json:"durationMs"`
	IsSuccess     bool   `json:"isSuccess"`
	ExitCode      int    `json:"exitCode"`
	OutputSummary string `json:"outputSummary,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Comments      string `json:"comments,omitempty"`
}

// StartupDBPath returns the full path to the startup SQLite DB.
func StartupDBPath() string {
	return ResolveSplitDbPath(SectionStartup, "default", "")
}

// OpenStartupSplitDB opens or initializes the startup split database.
func OpenStartupSplitDB() (*StartupSplitDB, error) {
	dbPath := StartupDBPath()
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open startup split db")
	}

	return initStartupSplitConn(conn, dbPath)
}

func initStartupSplitConn(conn *sql.DB, dbPath string) (*StartupSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	if err := createStartupSchema(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	db := &StartupSplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()

	return db, nil
}

func createStartupSchema(conn *sql.DB) error {
	if _, err := conn.Exec(constants.SQLCreateStartupItem); err != nil {
		return apperror.WrapSimple(err, "create StartupItem table")
	}

	if _, err := conn.Exec(constants.SQLCreateStartupLog); err != nil {
		return apperror.WrapSimple(err, "create StartupLog table")
	}

	return nil
}

func (db *StartupSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}
	defer master.Close()

	entry := SplitDatabaseEntry{
		DatabaseType:  "startup",
		DatabaseKey:   "startup_master",
		DatabasePath:  db.Path,
		Status:        "active",
		IsActive:      true,
		Description:   "Startup items and execution logs",
		SchemaVersion: 1,
	}
	_ = master.RegisterSplitDB(entry)
}

// Close terminates the database connection.
func (db *StartupSplitDB) Close() error {
	if db.conn == nil {
		return nil
	}

	return db.conn.Close()
}
