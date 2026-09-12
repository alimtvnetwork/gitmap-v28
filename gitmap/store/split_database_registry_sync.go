package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	_ "modernc.org/sqlite"
)

func (db *DB) dataDir() string {
	if db != nil && db.dbDir != "" {
		return db.dbDir
	}

	return BinaryDataDir()
}

// SyncKnownSplitDatabases discovers and updates metadata for all known split databases.
func (db *DB) SyncKnownSplitDatabases() error {
	if err := db.syncInstallationDB(); err != nil {
		return apperror.WrapSimple(err, "store.SyncKnownSplitDatabases.installation")
	}
	db.syncScheduleDBs()
	db.syncPipelineDBs()

	return nil
}

func (db *DB) syncInstallationDB() error {
	path := filepath.Join(db.dataDir(), installationDBFileName)
	instDB, err := OpenInstallationSplitDBAt(path)
	if err == nil && instDB != nil {
		_ = instDB.Close()
	}
	desc := "Dedicated system and developer tool installation split database"
	entry := inspectSplitDBFile("installation", "default", path, desc)

	return db.RegisterSplitDB(entry)
}

func (db *DB) syncScheduleDBs() {
	schedDir := filepath.Join(db.dataDir(), "schedules")
	entries, err := os.ReadDir(schedDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		db.syncSingleSchedule(schedDir, e)
	}
}

func (db *DB) syncSingleSchedule(dir string, e os.DirEntry) {
	if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
		return
	}
	slug := strings.TrimSuffix(e.Name(), ".db")
	path := filepath.Join(dir, e.Name())
	entry := inspectSplitDBFile("schedule", slug, path, "Schedule execution split database")
	_ = db.RegisterSplitDB(entry)
}

func (db *DB) syncPipelineDBs() {
	pipeDir := filepath.Join(db.dataDir(), "pipeline_db")
	entries, err := os.ReadDir(pipeDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		db.syncSinglePipeline(pipeDir, e)
	}
}

func (db *DB) syncSinglePipeline(dir string, e os.DirEntry) {
	if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
		return
	}
	name := strings.TrimSuffix(e.Name(), ".db")
	slug := strings.TrimPrefix(name, "pipeline_")
	path := filepath.Join(dir, e.Name())
	entry := inspectSplitDBFile("pipeline", slug, path, "Pipeline execution split database")
	_ = db.RegisterSplitDB(entry)
}

func inspectSplitDBFile(dbType, dbKey, path, desc string) SplitDatabaseEntry {
	entry := newDefaultSplitEntry(dbType, dbKey, path, desc)
	info, err := os.Stat(path)
	if err != nil {
		entry.Status = "offline"
		entry.IsActive = false

		return entry
	}
	entry.SizeBytes = info.Size()
	inspectDBInternals(path, &entry)

	return entry
}

func newDefaultSplitEntry(dbType, dbKey, path, desc string) SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType:  dbType,
		DatabaseKey:   dbKey,
		DatabasePath:  path,
		Description:   desc,
		Status:        "active",
		IsActive:      true,
		SchemaVersion: 1,
		LastSyncedAt:  time.Now().Unix(),
	}
}

func inspectDBInternals(path string, entry *SplitDatabaseEntry) {
	conn, err := OpenSQLiteDB(path)
	if err != nil {
		entry.Status = "error"
		entry.IsActive = false

		return
	}
	defer conn.Close()

	tables := getTableNames(conn)
	entry.TableCount = len(tables)
	entry.RecordCount = countRecordsForTables(conn, tables)
}

func getTableNames(conn *sql.DB) []string {
	q := "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';"
	rows, err := conn.Query(q)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return collectTableNames(rows)
}

func collectTableNames(rows *sql.Rows) []string {
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			names = append(names, name)
		}
	}

	return names
}

func countRecordsForTables(conn *sql.DB, tables []string) int64 {
	var total int64
	for _, tbl := range tables {
		total += countSingleTable(conn, tbl)
	}

	return total
}

func countSingleTable(conn *sql.DB, tbl string) int64 {
	var count int64
	q := fmt.Sprintf("SELECT COUNT(*) FROM [%s];", tbl)
	if err := conn.QueryRow(q).Scan(&count); err != nil {
		return 0
	}

	return count
}
