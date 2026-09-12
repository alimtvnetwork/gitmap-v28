package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

const (
	installationDBFileName = "installation.db"
	sqlCreateInstalledTool = `CREATE TABLE IF NOT EXISTS InstalledTool (
    InstalledToolId INTEGER PRIMARY KEY AUTOINCREMENT,
    Tool            TEXT NOT NULL UNIQUE,
    VersionMajor    INTEGER NOT NULL DEFAULT 0,
    VersionMinor    INTEGER NOT NULL DEFAULT 0,
    VersionPatch    INTEGER NOT NULL DEFAULT 0,
    VersionBuild    INTEGER NOT NULL DEFAULT 0,
    VersionString   TEXT NOT NULL DEFAULT '',
    PackageManager  TEXT NOT NULL DEFAULT '',
    InstallPath     TEXT NOT NULL DEFAULT '',
    Description     TEXT NULL,
    InstalledAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxInstalledTool_Tool ON InstalledTool(Tool);`

	sqlCreateInstallationLog = `CREATE TABLE IF NOT EXISTS InstallationLog (
    InstallationLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    Tool              TEXT NOT NULL,
    Action            TEXT NOT NULL,
    Version           TEXT NOT NULL DEFAULT '',
    PackageManager    TEXT NOT NULL DEFAULT '',
    DurationMs        INTEGER NOT NULL DEFAULT 0,
    IsSuccess         INTEGER NOT NULL DEFAULT 1,
    ExitCode          INTEGER NOT NULL DEFAULT 0,
    Stdout            TEXT NULL,
    Stderr            TEXT NULL,
    CommandLine       TEXT NULL,
    Notes             TEXT NULL,
    Comments          TEXT NULL,
    CreatedAt         TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxInstallationLog_Tool ON InstallationLog(Tool);
CREATE INDEX IF NOT EXISTS IdxInstallationLog_Action ON InstallationLog(Action);
CREATE INDEX IF NOT EXISTS IdxInstallationLog_CreatedAt ON InstallationLog(CreatedAt);`
)

// InstallationSplitDB wraps an isolated SQLite database connection for installations.
type InstallationSplitDB struct {
	conn *sql.DB
	Path string
}

// InstallationDBPath returns the full path to the split installation database.
func InstallationDBPath() string {
	return filepath.Join(BinaryDataDir(), installationDBFileName)
}

// OpenInstallationSplitDB opens the canonical split database for tool installations.
func OpenInstallationSplitDB() (*InstallationSplitDB, error) {
	return OpenInstallationSplitDBAt(InstallationDBPath())
}

// OpenInstallationSplitDBAt opens or creates a split database at a specific path.
func OpenInstallationSplitDBAt(dbPath string) (*InstallationSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.open")
	}

	return initInstallationSplitConn(conn, dbPath)
}

func initInstallationSplitConn(conn *sql.DB, dbPath string) (*InstallationSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "installation_split.config")
	}

	db := &InstallationSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return db, nil
}

// InitSchema creates the InstalledTool and InstallationLog tables if absent.
func (s *InstallationSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreateInstalledTool); err != nil {
		return apperror.WrapSimple(err, "installation_split.initInstalledTool")
	}

	if _, err := s.conn.Exec(sqlCreateInstallationLog); err != nil {
		return apperror.WrapSimple(err, "installation_split.initInstallationLog")
	}

	return nil
}

// Close closes the underlying split database connection.
func (s *InstallationSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Conn returns the raw database connection.
func (s *InstallationSplitDB) Conn() *sql.DB {
	return s.conn
}
