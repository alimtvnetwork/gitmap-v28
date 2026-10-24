package heatmap

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveRepoCacheDir resolves the local path to the repo-cache directory.
func ResolveRepoCacheDir() (string, error) {
	if cached, ok := resolveRepoCacheFromSpecialDB(); ok {
		return cached, nil
	}
	for _, c := range []string{`D:\work\repo-cache`, `d:\work\repo-cache`, `../repo-cache`} {
		if info, statErr := os.Stat(c); statErr == nil && info.IsDir() {
			return filepath.Abs(c)
		}
	}
	return filepath.Abs("../repo-cache")
}

// ExportPendingChangesToTempDB exports in-flight changes into a temporary snapshot database.
func ExportPendingChangesToTempDB(aiName, repoSlug string) (string, error) {
	repoCachePath, err := ResolveRepoCacheDir()
	if err != nil {
		return "", err
	}
	lastTs := "1970-01-01T00:00:00Z"
	if tr, _ := GetLocalTrackingStatus(repoSlug, aiName); tr != nil && tr.LastCommitTimestamp != "" {
		lastTs = tr.LastCommitTimestamp
	}
	aiFolder := ResolveAiFolderName(aiName)
	rootDb := filepath.Join(repoCachePath, ".heatmap", aiFolder, "root-heatmap.db")
	repoDb := filepath.Join(repoCachePath, ".heatmap", aiFolder, repoSlug, "heatmap.db")
	tempDb := filepath.Join(repoCachePath, ".heatmap", fmt.Sprintf("temp-sync-%d.db", time.Now().UnixNano()))
	if _, exportErr := ExportUncommittedChanges(rootDb, repoDb, tempDb, lastTs); exportErr != nil {
		return "", exportErr
	}
	return tempDb, nil
}

// ExportUncommittedChanges queries and exports pending deltas into tempDbPath.
func ExportUncommittedChanges(rootDbPath, repoDbPath, tempDbPath, lastCommitTimestamp string) (int, error) {
	_ = os.MkdirAll(filepath.Dir(tempDbPath), 0755)
	openTempRes := store.OpenSQLiteDB(tempDbPath)
	if openTempRes.IsFailure() {
		return 0, openTempRes.AppError()
	}
	tempConn := openTempRes.Value
	defer tempConn.Close()

	if err := createTempExportTables(tempConn); err != nil {
		return 0, err
	}
	rowsExported := exportParentTasks(rootDbPath, tempConn, lastCommitTimestamp) +
		exportActionLogs(repoDbPath, tempConn, lastCommitTimestamp)
	if err := verifyDbIntegrity(tempConn); err != nil {
		return 0, err
	}
	return rowsExported, nil
}

func createTempExportTables(conn *sql.DB) error {
	ddl := `CREATE TABLE IF NOT EXISTS HeatmapParentTask (TaskSlug TEXT PRIMARY KEY, TaskName TEXT, RunDirectory TEXT, RootDbPath TEXT, Status TEXT, Notes TEXT, Comments TEXT, CreatedAt TEXT, UpdatedAt TEXT);
	CREATE TABLE IF NOT EXISTS HeatmapActionLog (HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT, ParentLogId INTEGER, TaskSlug TEXT, RepoName TEXT, ActionType TEXT, FilePath TEXT, StartLine INTEGER, EndLine INTEGER, CodeSnippet TEXT, CommandText TEXT, ReasoningState TEXT, ReasoningNotes TEXT, ExpectedOutcome TEXT, ActualOutcome TEXT, IsExpectationMet INTEGER, ExecutionDepth INTEGER, ExecutionOrder INTEGER, Notes TEXT, Comments TEXT, CreatedAt TEXT);`
	_, err := conn.Exec(ddl)
	return err
}

func exportParentTasks(rootDbPath string, tempConn *sql.DB, timestamp string) int {
	openRootRes := store.OpenSQLiteDB(rootDbPath)
	if openRootRes.IsFailure() {
		return 0
	}
	rootConn := openRootRes.Value
	defer rootConn.Close()
	qRows, qErr := rootConn.Query(`SELECT TaskSlug, TaskName, RunDirectory, RootDbPath, Status, COALESCE(Notes,''), COALESCE(Comments,''), CreatedAt, UpdatedAt FROM HeatmapParentTask WHERE CreatedAt > ? OR UpdatedAt > ?`, timestamp, timestamp)
	if qErr != nil {
		return 0
	}
	defer qRows.Close()
	count := 0
	for qRows.Next() {
		var s, n, rd, rp, st, not, c, ca, ua string
		if qRows.Scan(&s, &n, &rd, &rp, &st, &not, &c, &ca, &ua) != nil {
			continue
		}
		if insertParentTaskRow(tempConn, s, n, rd, rp, st, not, c, ca, ua) {
			count++
		}
	}
	return count
}
func exportActionLogs(repoDbPath string, tempConn *sql.DB, timestamp string) int {
	openRepoRes := store.OpenSQLiteDB(repoDbPath)
	if openRepoRes.IsFailure() {
		return 0
	}
	repoConn := openRepoRes.Value
	defer repoConn.Close()
	qRows, qErr := repoConn.Query(`SELECT ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine, COALESCE(CodeSnippet,''), COALESCE(CommandText,''), ReasoningState, COALESCE(ReasoningNotes,''), COALESCE(ExpectedOutcome,''), COALESCE(ActualOutcome,''), IsExpectationMet, ExecutionDepth, ExecutionOrder, COALESCE(Notes,''), COALESCE(Comments,''), CreatedAt FROM HeatmapActionLog WHERE CreatedAt > ?`, timestamp)
	if qErr != nil {
		return 0
	}
	defer qRows.Close()
	return copyActionLogRows(qRows, tempConn)
}
func verifyDbIntegrity(conn *sql.DB) error {
	var checkResult string
	err := conn.QueryRow("PRAGMA integrity_check").Scan(&checkResult)
	if err != nil || checkResult != "ok" {
		return apperror.NewValidation("heatmap_conflict", "E_HEATMAP_TX_EXPORT", "temporary snapshot integrity failed: "+checkResult)
	}
	return nil
}

// DetectRemoteDivergence checks if remote changes exist in repoCachePath.
func DetectRemoteDivergence(repoCachePath, lastCommittedSha string) (bool, error) {
	fetch := exec.Command("git", "fetch", "origin")
	fetch.Dir = repoCachePath
	_ = fetch.Run()
	localSha, err := getGitRev(repoCachePath, "HEAD")
	if err != nil {
		return false, err
	}
	remoteSha, _ := getGitRev(repoCachePath, "origin/main")
	if remoteSha == "" {
		remoteSha, _ = getGitRev(repoCachePath, "origin/master")
	}
	return (remoteSha != "" && localSha != remoteSha) || (lastCommittedSha != "" && localSha != lastCommittedSha), nil
}

func getGitRev(dir, rev string) (string, error) {
	cmd := exec.Command("git", "rev-parse", rev)
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// DiscardUncommittedHeatmapDbs reverts modified database files in .heatmap/.
func DiscardUncommittedHeatmapDbs(repoCachePath string) error {
	cmd := exec.Command("git", "restore", ".heatmap/")
	cmd.Dir = repoCachePath
	if _, err := cmd.CombinedOutput(); err != nil {
		return discardHeatmapDbsFallback(repoCachePath)
	}
	return nil
}

// PullRebaseRepoCache fetches and rebases remote commits with retries.
func PullRebaseRepoCache(repoCachePath string) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		cmd := exec.Command("git", "pull", "--rebase", "origin", "main")
		cmd.Dir = repoCachePath
		if out, err := cmd.CombinedOutput(); err == nil {
			return nil
		} else {
			lastErr = apperror.WrapSimple(err, fmt.Sprintf("rebase failed: %s", string(out)))
		}
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}
	return apperror.WrapSimple(lastErr, "heatmap_conflict.pullRebase")
}

// ReconcileAndReplayChanges replays exported deltas from temp DB into target databases.
func ReconcileAndReplayChanges(tempDbPath, aiName, repoSlug string) error {
	repoCachePath, err := ResolveRepoCacheDir()
	if err != nil {
		return err
	}
	aiFolder := ResolveAiFolderName(aiName)
	rootDb := filepath.Join(repoCachePath, ".heatmap", aiFolder, "root-heatmap.db")
	repoDb := filepath.Join(repoCachePath, ".heatmap", aiFolder, repoSlug, "heatmap.db")
	return ReplayExportedTransactions(rootDb, repoDb, tempDbPath)
}

// ReplayExportedTransactions reads rows from tempDbPath and upserts them into destination databases.
func ReplayExportedTransactions(rootDbPath, repoDbPath, tempDbPath string) error {
	if _, err := os.Stat(tempDbPath); err != nil {
		return nil
	}
	openTempRes := store.OpenSQLiteDB(tempDbPath)
	if openTempRes.IsFailure() {
		return openTempRes.AppError()
	}
	tempConn := openTempRes.Value
	defer tempConn.Close()
	_ = replayParentTasks(tempConn, rootDbPath)
	_ = replayActionLogs(tempConn, repoDbPath)
	return nil
}

func replayParentTasks(tempConn *sql.DB, rootDbPath string) error {
	openRootRes := store.OpenSQLiteDB(rootDbPath)
	if openRootRes.IsFailure() {
		return openRootRes.AppError()
	}
	rootConn := openRootRes.Value
	defer rootConn.Close()
	rows, qErr := tempConn.Query(`SELECT TaskSlug, TaskName, RunDirectory, RootDbPath, Status, COALESCE(Notes,''), COALESCE(Comments,''), CreatedAt, UpdatedAt FROM HeatmapParentTask`)
	if qErr != nil {
		return qErr
	}
	defer rows.Close()
	upsertSQL := `INSERT INTO HeatmapParentTask (TaskSlug, TaskName, RunDirectory, RootDbPath, Status, Notes, Comments, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(TaskSlug) DO UPDATE SET Status = excluded.Status, Notes = excluded.Notes, Comments = excluded.Comments, UpdatedAt = excluded.UpdatedAt`
	for rows.Next() {
		var s, n, rd, rp, st, not, c, ca, ua string
		if scanErr := rows.Scan(&s, &n, &rd, &rp, &st, &not, &c, &ca, &ua); scanErr != nil {
			continue
		}
		if _, execErr := rootConn.Exec(upsertSQL, s, n, rd, rp, st, not, c, ca, ua); execErr != nil {
			return execErr
		}
	}
	return nil
}

func replayActionLogs(tempConn *sql.DB, repoDbPath string) error {
	openRepoRes := store.OpenSQLiteDB(repoDbPath)
	if openRepoRes.IsFailure() {
		return openRepoRes.AppError()
	}
	repoConn := openRepoRes.Value
	defer repoConn.Close()
	rows, qErr := tempConn.Query(`SELECT ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine, COALESCE(CodeSnippet,''), COALESCE(CommandText,''), ReasoningState, COALESCE(ReasoningNotes,''), COALESCE(ExpectedOutcome,''), COALESCE(ActualOutcome,''), IsExpectationMet, ExecutionDepth, ExecutionOrder, COALESCE(Notes,''), COALESCE(Comments,''), CreatedAt FROM HeatmapActionLog`)
	if qErr != nil {
		return qErr
	}
	defer rows.Close()
	_ = copyActionLogRows(rows, repoConn)
	return nil
}

// CleanupTempDB removes the ephemeral snapshot SQLite database and journal files.
func CleanupTempDB(tempDbPath string) error {
	if tempDbPath != "" {
		_ = os.Remove(tempDbPath)
		_ = os.Remove(tempDbPath + "-wal")
		_ = os.Remove(tempDbPath + "-shm")
	}
	return nil
}

// RunRcHeatmapSyncWithConflictResolution executes the full 6-stage reconciliation protocol.
func RunRcHeatmapSyncWithConflictResolution(aiName, repoSlug, taskSlug string) error {
	repoCachePath, err := ResolveRepoCacheDir()
	if err != nil {
		return err
	}
	tracking, err := GetLocalTrackingStatus(repoSlug, aiName)
	if err != nil {
		return err
	}
	if diverged, _ := DetectRemoteDivergence(repoCachePath, tracking.LastCommittedSha); diverged || tracking.IsDirty {
		_ = SetSyncStatus(repoSlug, aiName, "SYNCING", true)
		replayExportedChanges(repoCachePath, aiName, repoSlug)
	}
	taskSummary := taskSlug
	if taskSummary == "" {
		taskSummary = "sync decision telemetry"
	}
	if _, commitErr := CommitRcHeatmap(repoCachePath, repoSlug, aiName, taskSummary); commitErr != nil {
		_ = SetSyncStatus(repoSlug, aiName, "ERROR", true)
		return commitErr
	}
	return nil
}
