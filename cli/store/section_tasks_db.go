// Package store — section_tasks_db.go manages section-scoped tasks databases.
package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// SectionTasksDB wraps an isolated SQLite database connection for section tasks.
type SectionTasksDB struct {
	conn    *sql.DB
	Section string
	Path    string
}

// OpenSectionTasksDB opens or creates a section-scoped tasks database (.gitmap/data/<section>/<section>-tasks.db).
func OpenSectionTasksDB(section, repoRoot string) (*SectionTasksDB, error) {
	dbPath := ResolveSectionTasksDbPath(section, repoRoot)
	return OpenSectionTasksDBAt(section, dbPath)
}

// OpenSectionTasksDBAt opens or creates a section-scoped tasks database at a specific path.
func OpenSectionTasksDBAt(section, dbPath string) (*SectionTasksDB, error) {
	mkErr := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if mkErr != nil {
		return nil, apperror.WrapSimple(mkErr, "section_tasks.mkdir")
	}

	conn, openErr := sql.Open("sqlite", dbPath)
	if openErr != nil {
		return nil, apperror.WrapSimple(openErr, "section_tasks.open")
	}

	return initSectionTasksConn(conn, section, dbPath)
}

func initSectionTasksConn(conn *sql.DB, section, dbPath string) (*SectionTasksDB, error) {
	cfgErr := ConfigureSQLiteConn(conn)
	if cfgErr != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(cfgErr, "section_tasks.config")
	}

	db := &SectionTasksDB{conn: conn, Section: section, Path: dbPath}
	initErr := db.InitSchema()
	if initErr != nil {
		_ = conn.Close()
		return nil, initErr
	}

	return db, nil
}

// InitSchema creates standard task queue and task history tables for the section.
func (s *SectionTasksDB) InitSchema() error {
	resQueue := ExecWrapper(s.conn, sqlCreateTaskQueue)
	if resQueue.IsFailure {
		return apperror.WrapSimple(resQueue.Error, "section_tasks.initTaskQueue")
	}

	resHist := ExecWrapper(s.conn, sqlCreateTaskHistory)
	if resHist.IsFailure {
		return apperror.WrapSimple(resHist.Error, "section_tasks.initTaskHistory")
	}

	return nil
}

// Close closes the underlying database connection.
func (s *SectionTasksDB) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

// Conn returns the underlying database connection.
func (s *SectionTasksDB) Conn() *sql.DB {
	if s == nil {
		return nil
	}
	return s.conn
}
