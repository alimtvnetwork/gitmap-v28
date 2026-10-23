package cmdpull

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// PullTaskRecord represents an execution session row in PullTask.
type PullTaskRecord struct {
	PullTaskId   int64  `json:"pullTaskId"`
	TaskUuid     string `json:"taskUuid"`
	CommandType  string `json:"commandType"`
	TargetNode   string `json:"targetNode"`
	MachineId    string `json:"machineId"`
	Status       string `json:"status"`
	Pid          int    `json:"pid"`
	TotalRepos   int    `json:"totalRepos"`
	PulledRepos  int    `json:"pulledRepos"`
	EtaSeconds   int    `json:"etaSeconds"`
	StartedAt    string `json:"startedAt"`
	CompletedAt  string `json:"completedAt"`
	ErrorMessage string `json:"errorMessage"`
	HeartbeatAt  string `json:"heartbeatAt"`
	CreatedAt    string `json:"createdAt"`
}

// PullTaskRepoRecord represents a repository-level progress row in PullTaskRepo.
type PullTaskRepoRecord struct {
	PullTaskRepoId int64  `json:"pullTaskRepoId"`
	PullTaskId     int64  `json:"pullTaskId"`
	RepoSlug       string `json:"repoSlug"`
	RepoPath       string `json:"repoPath"`
	Status         string `json:"status"`
	CommitSha      string `json:"commitSha"`
	CommitTitle    string `json:"commitTitle"`
	FilesChanged   int    `json:"filesChanged"`
	DurationMs     int64  `json:"durationMs"`
	ErrorMessage   string `json:"errorMessage"`
	PulledAt       string `json:"pulledAt"`
	CreatedAt      string `json:"createdAt"`
}

// PullEtaHistoryRecord models a historical duration entry in PullEtaHistory.
type PullEtaHistoryRecord struct {
	PullEtaHistoryId int64   `json:"pullEtaHistoryId"`
	CommandType      string  `json:"commandType"`
	TargetNode       string  `json:"targetNode"`
	TotalRepos       int     `json:"totalRepos"`
	DurationSeconds  float64 `json:"durationSeconds"`
	CreatedAt        string  `json:"createdAt"`
}

// PullTaskSplitDB manages SQLite-backed task concurrency and progress in gitmap-pull.db.
type PullTaskSplitDB struct {
	conn *sql.DB
	path string
}

const (
	sqlCreatePullTask = `CREATE TABLE IF NOT EXISTS PullTask (
    PullTaskId      INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskUuid        TEXT NOT NULL UNIQUE,
    CommandType     TEXT NOT NULL,
    TargetNode      TEXT NOT NULL DEFAULT 'local',
    MachineId       TEXT NOT NULL DEFAULT '',
    Status          TEXT NOT NULL DEFAULT 'pending',
    Pid             INTEGER NOT NULL DEFAULT 0,
    TotalRepos      INTEGER NOT NULL DEFAULT 0,
    PulledRepos     INTEGER NOT NULL DEFAULT 0,
    EtaSeconds      INTEGER NOT NULL DEFAULT 0,
    StartedAt       TEXT NULL,
    CompletedAt     TEXT NULL,
    ErrorMessage    TEXT NULL,
    HeartbeatAt     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxPullTask_Status ON PullTask(Status);
CREATE INDEX IF NOT EXISTS IdxPullTask_Heartbeat ON PullTask(HeartbeatAt);
CREATE INDEX IF NOT EXISTS IdxPullTask_CreatedAt ON PullTask(CreatedAt);`

	sqlCreatePullTaskRepo = `CREATE TABLE IF NOT EXISTS PullTaskRepo (
    PullTaskRepoId  INTEGER PRIMARY KEY AUTOINCREMENT,
    PullTaskId      INTEGER NOT NULL,
    RepoSlug        TEXT NOT NULL,
    RepoPath        TEXT NOT NULL,
    Status          TEXT NOT NULL DEFAULT 'pending',
    CommitSha       TEXT NOT NULL DEFAULT '',
    CommitTitle     TEXT NOT NULL DEFAULT '',
    FilesChanged    INTEGER NOT NULL DEFAULT 0,
    DurationMs      INTEGER NOT NULL DEFAULT 0,
    ErrorMessage    TEXT NULL,
    PulledAt        TEXT NULL,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(PullTaskId) REFERENCES PullTask(PullTaskId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS IdxPullTaskRepo_PullTaskId ON PullTaskRepo(PullTaskId);
CREATE INDEX IF NOT EXISTS IdxPullTaskRepo_RepoSlug ON PullTaskRepo(RepoSlug);
CREATE INDEX IF NOT EXISTS IdxPullTaskRepo_Status ON PullTaskRepo(Status);`

	sqlCreatePullEtaHistory = `CREATE TABLE IF NOT EXISTS PullEtaHistory (
    PullEtaHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommandType      TEXT NOT NULL,
    TargetNode       TEXT NOT NULL DEFAULT 'local',
    TotalRepos       INTEGER NOT NULL,
    DurationSeconds  REAL NOT NULL,
    CreatedAt        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxPullEtaHistory_Lookup ON PullEtaHistory(CommandType, TargetNode, CreatedAt);`
)

// OpenPullTaskSplitDB opens the canonical split database for pull task coordination.
func OpenPullTaskSplitDB() (*PullTaskSplitDB, *apperror.AppError) {
	return OpenPullTaskSplitDBAt(store.PullDBPath())
}

// OpenPullTaskSplitDBAt opens or creates the pull task split database at a specific path.
func OpenPullTaskSplitDBAt(dbPath string) (*PullTaskSplitDB, *apperror.AppError) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.open")
	}

	return initPullTaskSplitConn(conn, dbPath)
}

func initPullTaskSplitConn(conn *sql.DB, dbPath string) (*PullTaskSplitDB, *apperror.AppError) {
	if err := configurePullTaskConn(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return setupPullTaskDB(conn, dbPath)
}

func setupPullTaskDB(conn *sql.DB, dbPath string) (*PullTaskSplitDB, *apperror.AppError) {
	db := &PullTaskSplitDB{conn: conn, path: dbPath}
	if err := db.InitSchema(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return db, nil
}

func configurePullTaskConn(conn *sql.DB) *apperror.AppError {
	conn.SetMaxOpenConns(1)

	pragmas := []string{
		constants.SQLPragmaJournalWAL,
		constants.SQLPragmaSynchronousNor,
		constants.SQLPragmaBusyTimeout10s,
		constants.SQLEnableFK,
	}

	return execPragmas(conn, pragmas)
}

func execPragmas(conn *sql.DB, pragmas []string) *apperror.AppError {
	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return apperror.WrapSimple(err, "pull_task_db.pragma")
		}
	}

	return nil
}

// InitSchema ensures PullTask, PullTaskRepo, and PullEtaHistory tables exist.
func (s *PullTaskSplitDB) InitSchema() *apperror.AppError {
	statements := []string{sqlCreatePullTask, sqlCreatePullTaskRepo, sqlCreatePullEtaHistory}
	for _, stmt := range statements {
		if _, err := s.conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "pull_task_db.schema")
		}
	}

	return nil
}

// Close closes the underlying SQLite database connection handle.
func (s *PullTaskSplitDB) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Conn returns the raw database connection handle.
func (s *PullTaskSplitDB) Conn() *sql.DB {
	return s.conn
}

// Path returns the physical file path of the database on disk.
func (s *PullTaskSplitDB) Path() string {
	return s.path
}
