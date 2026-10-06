package cmdpull

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
		_ = db.InsertPullError(rec)
		_ = LogPullErrorToFile(rec)

		return
	}

	_ = db.ClearPullErrorsForRepo(state.RepoName, state.RepoPath)
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
	isSuccess := state.ErrorMsg == "" && state.Step != PullStepTypeError
	successCount := 0
	failedCount := 0
	if isSuccess {
		successCount = 1
	} else {
		failedCount = 1
	}

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

func resolveCommitMetadata(repoPath string) (string, string) {
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--pretty=format:%an%x1f%s")
	out, err := cmd.Output()
	if err != nil {
		return "", ""
	}

	parts := strings.Split(string(out), "\x1f")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), ""
	}

	return "", ""
}

func resolveCommitTrace(repoPath, oldSHA, newSHA string) string {
	diffTrace := queryRangeCommitTrace(repoPath, oldSHA, newSHA)
	if len(diffTrace) > 0 {
		return diffTrace
	}

	return queryLatestCommitTrace(repoPath)
}

func queryRangeCommitTrace(repoPath, oldSHA, newSHA string) string {
	if oldSHA == "" || newSHA == "" || oldSHA == newSHA {
		return ""
	}
	cmd := exec.Command("git", "-C", repoPath, "log", oldSHA+".."+newSHA, "--oneline")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func queryLatestCommitTrace(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--oneline")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func parseFilesChanged(changes, repoPath, oldSHA, newSHA string) int {
	count := extractFilesFromChangesString(changes)
	if count > 0 {
		return count
	}

	return queryDiffStatFilesCount(repoPath, oldSHA, newSHA)
}

func extractFilesFromChangesString(changes string) int {
	start := strings.Index(changes, "(")
	end := strings.Index(changes, ")")
	if start == -1 || end <= start+1 {
		return parseFilesKeywordCount(changes)
	}

	num, err := strconv.Atoi(strings.TrimSpace(changes[start+1 : end]))
	if err != nil {
		return parseFilesKeywordCount(changes)
	}

	return num
}

func parseFilesKeywordCount(changes string) int {
	if !strings.Contains(changes, "file") {
		return 0
	}

	fields := strings.Fields(changes)
	for i, f := range fields {
		count, isMatch := tryParseFileField(fields, i, f)
		if isMatch {
			return count
		}
	}

	return 0
}

func tryParseFileField(fields []string, i int, f string) (int, bool) {
	if i <= 0 || !strings.HasPrefix(f, "file") {
		return 0, false
	}

	num, err := strconv.Atoi(fields[i-1])
	if err != nil {
		return 0, false
	}

	return num, true
}

func queryDiffStatFilesCount(repoPath, oldSHA, newSHA string) int {
	if oldSHA == "" || newSHA == "" || oldSHA == newSHA {
		return 0
	}

	cmd := exec.Command("git", "-C", repoPath, "diff", "--shortstat", oldSHA+".."+newSHA)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	return parseFilesKeywordCount(string(out))
}
