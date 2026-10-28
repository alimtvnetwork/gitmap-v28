package cmdpull

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RecordPullBatchSession writes pull session and repository telemetry to gitmap-pull.db.
func RecordPullBatchSession(telemetry PullSessionTelemetry, states []*PullRepoState) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "pull_db_sync.open_db")
	}
	defer db.Close()
	runID, err := db.InsertPullRun(telemetry.ConvertToStoreRunRecord())
	if err != nil {
		return err
	}
	syncPullFailuresToDB(db, states)

	return db.InsertPullRepoRuns(runID, buildRepoRunRecords(states))
}

func syncPullFailuresToDB(db *store.PullSplitDB, states []*PullRepoState) {
	if db == nil {
		return
	}
	_ = db.EnsurePullErrorsTable()
	nodeID := resolveLocalNodeIdentifier()
	now := time.Now().UTC()
	for _, state := range states {
		syncSingleStateToDB(db, state, nodeID, now)
	}
}

func syncSingleStateToDB(db *store.PullSplitDB, state *PullRepoState, nodeID string, now time.Time) {
	if state == nil {
		return
	}
	if isStateFailure(state) {
		rec := buildPullErrorRecordFromState(state, nodeID, now)
		if insErr := db.InsertPullError(rec); insErr != nil {
			_ = LogPullErrorToFile(rec)
		} else {
			_ = LogPullErrorToFile(rec)
		}

		return
	}

	if clearErr := db.ClearPullErrorsForRepo(state.RepoName, state.RepoPath); clearErr != nil {
		fmt.Fprintf(os.Stderr, "  %s⚠ Notice: error clearing pull status for %s: %v%s\n",
			constants.ColorYellow, state.RepoName, clearErr, constants.ColorReset)
	}
}

func buildPullErrorRecordFromState(state *PullRepoState, nodeID string, now time.Time) store.PullErrorRecord {
	return store.PullErrorRecord{
		ErrorID:        fmt.Sprintf("err-%d", time.Now().UnixNano()),
		RepoSlug:       state.RepoName,
		RepoPath:       state.RepoPath,
		NodeID:         nodeID,
		NodeVersion:    constants.Version,
		ErrorType:      string(state.Step),
		ErrorText:      state.ErrorMsg,
		StackTrace:     resolveStateStackTrace(state),
		RemediationCmd: ResolvePullRemediationHint(state),
		CreatedAt:      now,
	}
}

func resolveLocalNodeIdentifier() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "local-01"
	}
	return host
}

func resolveStateStackTrace(state *PullRepoState) string {
	if state == nil {
		return ""
	}
	if strings.Contains(state.ErrorMsg, "\n") {
		return strings.TrimSpace(state.ErrorMsg)
	}
	if state.ErrorMsg != "" {
		return fmt.Sprintf("failed at step [%s]: %s", state.Step, state.ErrorMsg)
	}
	return fmt.Sprintf("failed at step [%s]", state.Step)
}

func isStateFailure(state *PullRepoState) bool {
	if state == nil {
		return false
	}
	return state.ErrorMsg != "" || state.Step == PullStepTypeError
}

// RecordSinglePullSession records a single CWD pull invocation to gitmap-pull.db.
func RecordSinglePullSession(cwd string, state *PullRepoState, cmdType string) error {
	telemetry := buildSingleSessionTelemetry(cwd, state, cmdType)

	return RecordPullBatchSession(telemetry, []*PullRepoState{state})
}

func buildSingleSessionTelemetry(cwd string, state *PullRepoState, cmdType string) PullSessionTelemetry {
	successCount, failedCount := resolveSingleSessionCounts(state)

	return PullSessionTelemetry{
		CommandType:   cmdType,
		WorkingDir:    cwd,
		TotalRepos:    1,
		PulledRepos:   1,
		SkippedRepos:  0,
		SuccessCount:  successCount,
		FailedCount:   failedCount,
		IsEfficient:   false,
		Duration:      state.Duration,
		GitMapVersion: constants.Version,
	}
}

func resolveSingleSessionCounts(state *PullRepoState) (int, int) {
	if state.ErrorMsg == "" && state.Step != PullStepTypeError {
		return 1, 0
	}

	return 0, 1
}

func buildRepoRunRecords(states []*PullRepoState) []store.PullRepoRunRecord {
	records := make([]store.PullRepoRunRecord, 0, len(states))
	for _, state := range states {
		records = append(records, convertStateToRepoRun(state))
	}

	return records
}

func convertStateToRepoRun(state *PullRepoState) store.PullRepoRunRecord {
	lastSHA := resolveStateLastSHA(state)
	filesChanged := parseFilesChanged(state.Changes, state.RepoPath, state.OldSHA, state.NewSHA)
	hasChanges := state.OldSHA != state.NewSHA || filesChanged > 0
	author, msg := resolveCommitMetadata(state.RepoPath)
	trace := resolveCommitTrace(state.RepoPath, state.OldSHA, state.NewSHA)

	return store.PullRepoRunRecord{
		RepoPath:          state.RepoPath,
		RepoName:          state.RepoName,
		PullStatus:        resolveStatePullStatus(state),
		FilesChanged:      filesChanged,
		LastCommitSha:     lastSHA,
		PreviousCommitSha: state.OldSHA,
		CommitMessage:     msg,
		CommitAuthor:      author,
		IsActive:          true,
		HasChanges:        hasChanges,
		DurationMs:        state.Duration.Milliseconds(),
		ErrorMessage:      state.ErrorMsg,
		Notes:             trace,
	}
}

func resolveStateLastSHA(state *PullRepoState) string {
	if state.NewSHA != "" {
		return state.NewSHA
	}

	return gitutil.GetLastCommitSHA(state.RepoPath)
}

func resolveStatePullStatus(state *PullRepoState) string {
	if state.ErrorMsg != "" || state.Step == PullStepTypeError {
		return "failed"
	}
	if state.IsDirty {
		return "skipped-dirty"
	}
	if state.Step == PullStepTypeSkipped {
		return "skipped"
	}
	if state.OldSHA != "" && state.OldSHA == state.NewSHA {
		return "up-to-date"
	}

	return "success"
}
