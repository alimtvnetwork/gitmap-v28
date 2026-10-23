package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// CommitHistorySplitDB manages high-performance storage for commit records and touches.
type CommitHistorySplitDB struct {
	conn *sql.DB
	mu   sync.RWMutex
	Path string
}

// CommitHistoryDB is a type alias for CommitHistorySplitDB.
type CommitHistoryDB = CommitHistorySplitDB

// OpenCommitHistorySplitDB opens or initializes the commit history split database.
func OpenCommitHistorySplitDB(repoRoot string) (*CommitHistorySplitDB, error) {
	dbPath := ResolveCommitHistoryDbPath(repoRoot)

	return OpenCommitHistorySplitDBAt(dbPath)
}

// OpenCommitHistoryDB is a compatibility alias for OpenCommitHistorySplitDB.
func OpenCommitHistoryDB(repoRoot string) (*CommitHistorySplitDB, error) {
	return OpenCommitHistorySplitDB(repoRoot)
}

// OpenCommitHistorySplitDBAt opens the database at a specific filesystem path.
func OpenCommitHistorySplitDBAt(dbPath string) (*CommitHistorySplitDB, error) {
	conn, err := openConfiguredCommitConn(dbPath)
	if err != nil {
		return nil, err
	}

	db := &CommitHistorySplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()

	return db, nil
}

// OpenCommitHistoryDBAt is a compatibility alias for OpenCommitHistorySplitDBAt.
func OpenCommitHistoryDBAt(dbPath string) (*CommitHistorySplitDB, error) {
	return OpenCommitHistorySplitDBAt(dbPath)
}

func openConfiguredCommitConn(dbPath string) (*sql.DB, error) {
	conn, err := openCommitConn(dbPath)
	if err != nil {
		return nil, err
	}

	if errSetup := setupCommitHistoryDB(conn); errSetup != nil {
		_ = conn.Close()

		return nil, errSetup
	}

	return conn, nil
}

func openCommitConn(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "store.openCommitConn.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.openCommitConn.open")
	}

	conn.SetMaxOpenConns(1)

	return conn, nil
}

func setupCommitHistoryDB(conn *sql.DB) error {
	if err := applyCommitHistoryPragmas(conn); err != nil {
		return err
	}

	return initCommitHistorySchema(conn)
}

var commitHistoryPragmas = []string{
	"PRAGMA journal_mode = WAL;",
	"PRAGMA synchronous = NORMAL;",
	"PRAGMA busy_timeout = 5000;",
	"PRAGMA mmap_size = 268435456;",
	"PRAGMA cache_size = -64000;",
	"PRAGMA foreign_keys = ON;",
}

func applyCommitHistoryPragmas(conn *sql.DB) error {
	for _, p := range commitHistoryPragmas {
		if _, err := conn.Exec(p); err != nil {
			return apperror.WrapSimple(err, "store.applyCommitHistoryPragmas")
		}
	}

	return nil
}

func (db *CommitHistorySplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}

	defer master.Close()

	entry := makeCommitHistoryRegistryEntry(db.Path)
	_ = master.RegisterSplitDB(entry)
}

func makeCommitHistoryRegistryEntry(dbPath string) SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType: "repodb",
		DatabaseKey:  "commit_history",
		DatabasePath: dbPath,
		Description:  "Repository commit history, file touches, and AI agent provenance split database",
		Status:       "active",
		IsActive:     true,
	}
}

// Close gracefully releases the database connection.
func (db *CommitHistorySplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return db.conn.Close()
}
