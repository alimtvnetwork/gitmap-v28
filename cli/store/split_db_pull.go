package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	PullDBFileName = "gitmap-pull.db"

	sqlCreatePullRun = `CREATE TABLE IF NOT EXISTS PullRun (
    PullRunId       INTEGER PRIMARY KEY AUTOINCREMENT,
    CommandType     TEXT NOT NULL,
    WorkingDir      TEXT NOT NULL,
    TotalRepos      INTEGER NOT NULL DEFAULT 0,
    PulledRepos     INTEGER NOT NULL DEFAULT 0,
    SkippedRepos    INTEGER NOT NULL DEFAULT 0,
    SuccessCount    INTEGER NOT NULL DEFAULT 0,
    FailedCount     INTEGER NOT NULL DEFAULT 0,
    IsEfficient     INTEGER NOT NULL DEFAULT 0,
    DurationMs      INTEGER NOT NULL DEFAULT 0,
    GitMapVersion   TEXT NOT NULL DEFAULT '',
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxPullRun_CreatedAt ON PullRun(CreatedAt);
CREATE INDEX IF NOT EXISTS IdxPullRun_CommandType ON PullRun(CommandType);`

	sqlCreatePullRepoRun = `CREATE TABLE IF NOT EXISTS PullRepoRun (
    PullRepoRunId     INTEGER PRIMARY KEY AUTOINCREMENT,
    PullRunId         INTEGER NOT NULL,
    RepoPath          TEXT NOT NULL,
    RepoName          TEXT NOT NULL,
    PullStatus        TEXT NOT NULL DEFAULT 'success',
    FilesChanged      INTEGER NOT NULL DEFAULT 0,
    LastCommitSha     TEXT NOT NULL DEFAULT '',
    PreviousCommitSha TEXT NOT NULL DEFAULT '',
    CommitMessage     TEXT NOT NULL DEFAULT '',
    CommitAuthor      TEXT NOT NULL DEFAULT '',
    IsActive          INTEGER NOT NULL DEFAULT 1,
    HasChanges        INTEGER NOT NULL DEFAULT 0,
    DurationMs        INTEGER NOT NULL DEFAULT 0,
    ErrorMessage      TEXT NULL,
    Notes             TEXT NULL,
    Comments          TEXT NULL,
    CreatedAt         TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(PullRunId) REFERENCES PullRun(PullRunId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxPullRepoRun_PullRunId ON PullRepoRun(PullRunId);
CREATE INDEX IF NOT EXISTS IdxPullRepoRun_RepoPath ON PullRepoRun(RepoPath);
CREATE INDEX IF NOT EXISTS IdxPullRepoRun_CreatedAt ON PullRepoRun(CreatedAt);
CREATE INDEX IF NOT EXISTS IdxPullRepoRun_HasChanges ON PullRepoRun(HasChanges);`
)

// PullDBPath returns the canonical file path to gitmap-pull.db in the data folder.
func PullDBPath() string {
	dir := BinaryDataDir()
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, PullDBFileName)
}

// OpenPullSplitDB opens the canonical split database for git pull telemetry.
func OpenPullSplitDB() (*PullSplitDB, error) {
	return OpenPullSplitDBAt(PullDBPath())
}

// OpenPullSplitDBAt opens or creates a pull split database at a specific path.
func OpenPullSplitDBAt(dbPath string) (*PullSplitDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.open")
	}

	return initPullSplitConn(conn, dbPath)
}

func initPullSplitConn(conn *sql.DB, dbPath string) (*PullSplitDB, error) {
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, apperror.WrapSimple(err, "pull_split.config")
	}

	db := &PullSplitDB{conn: conn, Path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return db, nil
}

// InitSchema creates the PullRun and PullRepoRun tables if absent.
func (s *PullSplitDB) InitSchema() error {
	if _, err := s.conn.Exec(sqlCreatePullRun); err != nil {
		return apperror.WrapSimple(err, "pull_split.init_pull_run")
	}

	if _, err := s.conn.Exec(sqlCreatePullRepoRun); err != nil {
		return apperror.WrapSimple(err, "pull_split.init_pull_repo_run")
	}

	return nil
}

// Close closes the underlying SQLite database connection.
func (s *PullSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Conn returns the raw database connection handle.
func (s *PullSplitDB) Conn() *sql.DB {
	return s.conn
}
