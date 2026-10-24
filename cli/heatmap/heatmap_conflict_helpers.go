package heatmap

import (
	"database/sql"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// resolveRepoCacheFromSpecialDB returns the repo-cache path recorded in the
// special-repos split DB, reporting false when no usable directory is recorded.
func resolveRepoCacheFromSpecialDB() (string, bool) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return "", false
	}
	defer db.Close()
	rec, getErr := db.GetSpecialRepo("rc")
	if getErr != nil || rec.LocalPath == "" {
		return "", false
	}
	info, statErr := os.Stat(rec.LocalPath)
	if statErr != nil || !info.IsDir() {
		return "", false
	}
	return rec.LocalPath, true
}

// insertParentTaskRow inserts one exported parent-task row, reporting whether
// the insert succeeded.
func insertParentTaskRow(tempConn *sql.DB, s, n, rd, rp, st, not, c, ca, ua string) bool {
	_, execErr := tempConn.Exec(`INSERT OR REPLACE INTO HeatmapParentTask VALUES (?,?,?,?,?,?,?,?,?)`, s, n, rd, rp, st, not, c, ca, ua)
	return execErr == nil
}

// discardHeatmapDbsFallback restores .heatmap/ via checkout when restore fails.
func discardHeatmapDbsFallback(repoCachePath string) error {
	fbCmd := exec.Command("git", "checkout", "--", ".heatmap/")
	fbCmd.Dir = repoCachePath
	out, fbErr := fbCmd.CombinedOutput()
	if fbErr != nil {
		return apperror.WrapSimple(fbErr, "heatmap_conflict.discard: "+string(out))
	}
	return nil
}

// replayExportedChanges replays a pending-changes export into the target
// databases, doing nothing when there is nothing to export.
func replayExportedChanges(repoCachePath, aiName, repoSlug string) {
	tempDbPath, exportErr := ExportPendingChangesToTempDB(aiName, repoSlug)
	if exportErr != nil || tempDbPath == "" {
		return
	}
	_ = DiscardUncommittedHeatmapDbs(repoCachePath)
	_ = PullRebaseRepoCache(repoCachePath)
	_ = ReconcileAndReplayChanges(tempDbPath, aiName, repoSlug)
	_ = CleanupTempDB(tempDbPath)
}

func copyActionLogRows(rows *sql.Rows, destConn *sql.DB) int {
	insertSQL := `INSERT INTO HeatmapActionLog (ParentLogId, TaskSlug, RepoName, ActionType, FilePath, StartLine, EndLine, CodeSnippet, CommandText, ReasoningState, ReasoningNotes, ExpectedOutcome, ActualOutcome, IsExpectationMet, ExecutionDepth, ExecutionOrder, Notes, Comments, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	count := 0
	for rows.Next() {
		var pid sql.NullInt64
		var ts, rn, act, fp, cs, ct, rs, rnot, eo, ao, not, com, ca string
		var sl, el, isMet, ed, eoOrder int
		if rows.Scan(&pid, &ts, &rn, &act, &fp, &sl, &el, &cs, &ct, &rs, &rnot, &eo, &ao, &isMet, &ed, &eoOrder, &not, &com, &ca) != nil {
			continue
		}
		if _, execErr := destConn.Exec(insertSQL, pid, ts, rn, act, fp, sl, el, cs, ct, rs, rnot, eo, ao, isMet, ed, eoOrder, not, com, ca); execErr == nil {
			count++
		}
	}
	return count
}
