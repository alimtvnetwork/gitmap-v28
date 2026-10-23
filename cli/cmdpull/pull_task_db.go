package cmdpull

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
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

	sqlSelectActiveLockWithElapsed = `SELECT PullTaskId, TaskUuid, CommandType, TargetNode, MachineId, Status,
    Pid, TotalRepos, PulledRepos, EtaSeconds,
    COALESCE(StartedAt, ''), COALESCE(CompletedAt, ''),
    COALESCE(ErrorMessage, ''), HeartbeatAt, CreatedAt,
    (strftime('%s', 'now') - strftime('%s', HeartbeatAt)) as ElapsedSec
FROM PullTask
WHERE TargetNode = ? AND Status IN ('pending', 'running')
ORDER BY PullTaskId DESC LIMIT 1;`

	sqlInsertPullTask = `INSERT INTO PullTask (
    TaskUuid, CommandType, TargetNode, MachineId, Status, Pid,
    TotalRepos, PulledRepos, EtaSeconds, StartedAt, HeartbeatAt, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	sqlInsertPullTaskRepo = `INSERT INTO PullTaskRepo (
    PullTaskId, RepoSlug, RepoPath, Status, CommitSha, CommitTitle,
    FilesChanged, DurationMs, ErrorMessage, PulledAt, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	sqlSelectListRecentRepos = `SELECT PullTaskRepoId, PullTaskId, RepoSlug, RepoPath, Status,
    CommitSha, CommitTitle, FilesChanged, DurationMs,
    COALESCE(ErrorMessage, ''), COALESCE(PulledAt, ''), CreatedAt
FROM PullTaskRepo WHERE PullTaskId = ?
ORDER BY PullTaskRepoId DESC LIMIT ?;`

	sqlDeleteExpiredRepos = `DELETE FROM PullTaskRepo WHERE PullTaskId IN (
    SELECT PullTaskId FROM PullTask WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days')
);`

	sqlDeleteExpiredTasks = `DELETE FROM PullTask WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days');`

	sqlDeleteExpiredEta = `DELETE FROM PullEtaHistory WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days');`
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

// FindActiveLock queries for an active lock on targetNode with fresh heartbeat.
func (s *PullTaskSplitDB) FindActiveLock(targetNode string) (*PullTaskRecord, bool, *apperror.AppError) {
	node := resolveTargetNode(targetNode)
	row := s.conn.QueryRow(sqlSelectActiveLockWithElapsed, node)

	rec, elapsedSec, err := scanActiveLockRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, apperror.WrapSimple(err, "pull_task_db.find_lock")
	}

	return evaluateLockFreshness(s, rec, elapsedSec)
}

func evaluateLockFreshness(s *PullTaskSplitDB, rec *PullTaskRecord, elapsedSec int64) (*PullTaskRecord, bool, *apperror.AppError) {
	isFresh := elapsedSec < int64(constants.PullTaskStaleLockThresholdSec)
	if isFresh {
		return rec, true, nil
	}

	if err := s.StealStaleLock(rec.PullTaskId); err != nil {
		return nil, false, err
	}

	return nil, false, nil
}

// StealStaleLock marks an orphaned or crashed running task as failed.
func (s *PullTaskSplitDB) StealStaleLock(taskID int64) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, ErrorMessage = ?, CompletedAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`
	errMsg := "process terminated unexpectedly (stale heartbeat)"

	if _, err := s.conn.Exec(query, constants.PullTaskStatusFailed, errMsg, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.steal_stale_lock")
	}

	return nil
}

// TryAcquirePullTaskLock attempts to acquire the execution lock for a node.
func (s *PullTaskSplitDB) TryAcquirePullTaskLock(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, bool, *apperror.AppError) {
	activeRec, hasLock, err := s.FindActiveLock(targetNode)
	if err != nil {
		return nil, false, err
	}

	if hasLock {
		return activeRec, false, nil
	}

	return s.acquireFreshLock(cmdType, targetNode, totalRepos, etaSec)
}

func (s *PullTaskSplitDB) acquireFreshLock(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, bool, *apperror.AppError) {
	newRec, err := s.insertNewPullTask(cmdType, targetNode, totalRepos, etaSec)
	if err != nil {
		return nil, false, err
	}

	return newRec, true, nil
}

func (s *PullTaskSplitDB) insertNewPullTask(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, *apperror.AppError) {
	node := resolveTargetNode(targetNode)
	machineId := resolveMachineId()
	uuidStr := uuidutil.GenerateUUID()

	res, err := s.conn.Exec(sqlInsertPullTask, uuidStr, cmdType, node, machineId,
		constants.PullTaskStatusRunning, os.Getpid(), totalRepos, etaSec)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.insert_task")
	}

	return s.fetchInsertedTask(res)
}

func (s *PullTaskSplitDB) fetchInsertedTask(res sql.Result) (*PullTaskRecord, *apperror.AppError) {
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.last_insert_id")
	}

	return s.GetPullTask(id)
}

// UpdateTaskHeartbeat refreshes HeartbeatAt timestamp to maintain execution lease.
func (s *PullTaskSplitDB) UpdateTaskHeartbeat(taskID int64) *apperror.AppError {
	query := `UPDATE PullTask SET HeartbeatAt = CURRENT_TIMESTAMP WHERE PullTaskId = ?`
	if _, err := s.conn.Exec(query, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.heartbeat")
	}

	return nil
}

// UpdateTaskProgress updates pulled count, remaining ETA, and refreshes heartbeat.
func (s *PullTaskSplitDB) UpdateTaskProgress(taskID int64, pulledRepos, etaSec int) *apperror.AppError {
	query := `UPDATE PullTask
		SET PulledRepos = ?, EtaSeconds = ?, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, pulledRepos, etaSec, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.progress")
	}

	return nil
}

// RecordRepoPullResult inserts progress for a repository under an active task.
func (s *PullTaskSplitDB) RecordRepoPullResult(rec PullTaskRepoRecord) *apperror.AppError {
	title := sanitizeCommitTitle(rec.CommitTitle)
	sha := truncateSha(rec.CommitSha)

	_, err := s.conn.Exec(sqlInsertPullTaskRepo, rec.PullTaskId, rec.RepoSlug, rec.RepoPath, rec.Status,
		sha, title, rec.FilesChanged, rec.DurationMs, rec.ErrorMessage)
	if err != nil {
		return apperror.WrapSimple(err, "pull_task_db.record_repo")
	}

	return nil
}

// CompletePullTask marks the task completed and appends runtime to PullEtaHistory.
func (s *PullTaskSplitDB) CompletePullTask(taskID int64, totalRepos int, durationMs int64, cmdType, targetNode string) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, PulledRepos = TotalRepos, CompletedAt = CURRENT_TIMESTAMP, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, constants.PullTaskStatusCompleted, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.complete")
	}

	durationSec := float64(durationMs) / 1000.0

	return s.insertEtaHistory(cmdType, targetNode, totalRepos, durationSec)
}

func (s *PullTaskSplitDB) insertEtaHistory(cmdType, targetNode string, totalRepos int, durationSec float64) *apperror.AppError {
	node := resolveTargetNode(targetNode)
	query := `INSERT INTO PullEtaHistory (CommandType, TargetNode, TotalRepos, DurationSeconds, CreatedAt)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	if _, err := s.conn.Exec(query, cmdType, node, totalRepos, durationSec); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.insert_eta_history")
	}

	return nil
}

// FailPullTask marks the task failed with an explanatory error message.
func (s *PullTaskSplitDB) FailPullTask(taskID int64, errMsg string) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, ErrorMessage = ?, CompletedAt = CURRENT_TIMESTAMP, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, constants.PullTaskStatusFailed, errMsg, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.fail")
	}

	return nil
}

// GetPullTask fetches a single PullTaskRecord by its primary key.
func (s *PullTaskSplitDB) GetPullTask(taskID int64) (*PullTaskRecord, *apperror.AppError) {
	query := `SELECT PullTaskId, TaskUuid, CommandType, TargetNode, MachineId, Status,
		Pid, TotalRepos, PulledRepos, EtaSeconds,
		COALESCE(StartedAt, ''), COALESCE(CompletedAt, ''),
		COALESCE(ErrorMessage, ''), HeartbeatAt, CreatedAt
		FROM PullTask WHERE PullTaskId = ?`

	row := s.conn.QueryRow(query, taskID)
	rec, err := scanPullTaskRow(row)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.get_task")
	}

	return rec, nil
}

// ListRecentPullTaskRepos returns the most recent repository progress entries.
func (s *PullTaskSplitDB) ListRecentPullTaskRepos(taskID int64, limit int) ([]PullTaskRepoRecord, *apperror.AppError) {
	rows, err := s.conn.Query(sqlSelectListRecentRepos, taskID, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.list_repos")
	}

	defer rows.Close()

	return scanPullTaskRepoRows(rows)
}

// PruneExpiredPullHistory purges records older than retentionDays.
func (s *PullTaskSplitDB) PruneExpiredPullHistory(retentionDays int) (int64, *apperror.AppError) {
	daysStr := strconv.Itoa(resolveRetentionDays(retentionDays))

	deletedCount, err := s.purgeRetentionTables(daysStr)
	if err != nil {
		return 0, err
	}

	s.vacuumIfLargePurge(deletedCount)

	return deletedCount, nil
}

func (s *PullTaskSplitDB) purgeRetentionTables(daysStr string) (int64, *apperror.AppError) {
	queries := []string{sqlDeleteExpiredRepos, sqlDeleteExpiredTasks, sqlDeleteExpiredEta}
	var total int64

	for _, q := range queries {
		cnt, err := execDeleteQuery(s.conn, q, daysStr)
		if err != nil {
			return 0, err
		}

		total += cnt
	}

	return total, nil
}

func (s *PullTaskSplitDB) vacuumIfLargePurge(totalDeleted int64) {
	if totalDeleted > 500 {
		_, _ = s.conn.Exec("PRAGMA incremental_vacuum(100);")
	}
}

func execDeleteQuery(conn *sql.DB, query, arg string) (int64, *apperror.AppError) {
	res, err := conn.Exec(query, arg)
	if err != nil {
		return 0, apperror.WrapSimple(err, "pull_task_db.prune")
	}

	rows, _ := res.RowsAffected()

	return rows, nil
}

func scanActiveLockRow(row *sql.Row) (*PullTaskRecord, int64, error) {
	var t PullTaskRecord
	var elapsed sql.NullInt64

	err := row.Scan(&t.PullTaskId, &t.TaskUuid, &t.CommandType, &t.TargetNode, &t.MachineId, &t.Status,
		&t.Pid, &t.TotalRepos, &t.PulledRepos, &t.EtaSeconds,
		&t.StartedAt, &t.CompletedAt, &t.ErrorMessage, &t.HeartbeatAt, &t.CreatedAt, &elapsed)
	if err != nil {
		return nil, 0, err
	}

	return &t, elapsed.Int64, nil
}

func scanPullTaskRow(row *sql.Row) (*PullTaskRecord, error) {
	var t PullTaskRecord

	err := row.Scan(&t.PullTaskId, &t.TaskUuid, &t.CommandType, &t.TargetNode, &t.MachineId, &t.Status,
		&t.Pid, &t.TotalRepos, &t.PulledRepos, &t.EtaSeconds,
		&t.StartedAt, &t.CompletedAt, &t.ErrorMessage, &t.HeartbeatAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func scanPullTaskRepoRows(rows *sql.Rows) ([]PullTaskRepoRecord, *apperror.AppError) {
	var records []PullTaskRepoRecord
	for rows.Next() {
		r, err := scanSingleRepoRow(rows)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

func scanSingleRepoRow(rows *sql.Rows) (PullTaskRepoRecord, *apperror.AppError) {
	var r PullTaskRepoRecord
	err := rows.Scan(&r.PullTaskRepoId, &r.PullTaskId, &r.RepoSlug, &r.RepoPath, &r.Status,
		&r.CommitSha, &r.CommitTitle, &r.FilesChanged, &r.DurationMs,
		&r.ErrorMessage, &r.PulledAt, &r.CreatedAt)
	if err != nil {
		return r, apperror.WrapSimple(err, "pull_task_db.scan_repo")
	}

	return r, nil
}

func sanitizeCommitTitle(title string) string {
	clean := strings.TrimSpace(title)
	if len(clean) > 120 {
		return clean[:120]
	}

	return clean
}

func truncateSha(sha string) string {
	clean := strings.TrimSpace(sha)
	if len(clean) > 40 {
		return clean[:40]
	}

	return clean
}

func resolveTargetNode(node string) string {
	clean := strings.TrimSpace(node)
	if clean == "" {
		return "local"
	}

	return clean
}

func resolveMachineId() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return host
}

func resolveRetentionDays(days int) int {
	if days <= 0 {
		return constants.DefaultPullRetentionDays
	}

	return days
}
