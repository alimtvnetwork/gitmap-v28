package cmdpull

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunBatchPullFix loads last recorded pull failures from SQLite and runs remediation.
func RunBatchPullFix(args []string) error {
	failures, err := loadRecentPullFailuresFromDB()
	if err != nil {
		return err
	}
	hasFailures := len(failures) > 0
	if !hasFailures {
		fmt.Printf("  %sℹ No recent pull failures found to remediate.%s\n",
			constants.ColorCyan, constants.ColorReset)
		return nil
	}
	fmt.Printf("  %sFound %d failed repository(ies) from recent runs.%s\n",
		constants.ColorCyan, len(failures), constants.ColorReset)

	return dispatchBatchFixArgs(args, failures)
}

func loadRecentPullFailuresFromDB() ([]PullFailureSummary, *apperror.AppError) {
	db, errDB := store.OpenPullSplitDB()
	if errDB != nil {
		return nil, apperror.WrapSimple(errDB, "open_pull_db")
	}
	defer db.Close()
	recordsRes := db.QueryAllLatestPullErrors(50)
	if recordsRes.IsFailure() {
		return nil, recordsRes.AppError()
	}

	return convertRecordsToFailures(recordsRes.Items()), nil
}

func convertRecordsToFailures(records []store.PullErrorRecord) []PullFailureSummary {
	var failures []PullFailureSummary
	seen := make(map[string]bool, len(records))
	for _, r := range records {
		hasPath := r.RepoPath != ""
		if hasPath && !seen[r.RepoPath] {
			seen[r.RepoPath] = true
			repoID, remoteURL, _ := LookupRepoRemoteAndID(r.RepoPath)
			failures = append(failures, PullFailureSummary{
				RepoID:    repoID,
				RepoName:  r.RepoSlug,
				RepoPath:  r.RepoPath,
				RemoteURL: remoteURL,
				ErrorText: r.ErrorText,
				ErrorType: r.ErrorType,
			})
		}
	}

	return failures
}

func dispatchBatchFixArgs(args []string, failures []PullFailureSummary) error {
	if hasBatchFlag(args, "--rebase", "-r") {
		return RunBatchRemediateRebase(failures)
	}
	if hasBatchFlag(args, "--all", "-a", "--yes", "-y") {
		return RunBatchRemediateAll(failures)
	}
	if hasBatchFlag(args, "--step", "-s", "--one-by-one") {
		return RunStepByStepRemediate(failures)
	}

	return PromptInteractiveBatchFix(failures)
}

func hasBatchFlag(args []string, targets ...string) bool {
	for _, arg := range args {
		low := strings.ToLower(strings.TrimSpace(arg))
		for _, target := range targets {
			if low == target {
				return true
			}
		}
	}

	return false
}

// RemediateRepoByQuery resolves and remediates a single repository by name or path (including credential/SSH auto-fix).
func RemediateRepoByQuery(repoQuery string) (bool, error) {
	clean := strings.TrimSpace(repoQuery)
	if clean == "" {
		return false, nil
	}
	failures, _ := loadRecentPullFailuresFromDB()
	for _, f := range failures {
		handled, err := tryRemediateFailureMatch(f, clean)
		if handled {
			return true, err
		}
	}
	records := resolveAllTrackedRecords()
	matched := FindBySlug(records, clean)
	if len(matched) != 1 {
		return false, nil
	}
	return tryRemediateMatchedRecord(matched[0])
}

func tryRemediateFailureMatch(f PullFailureSummary, clean string) (bool, error) {
	if !strings.EqualFold(f.RepoName, clean) && !strings.EqualFold(f.RepoPath, clean) {
		return false, nil
	}
	if remediateSingleRepo(f) {
		return true, nil
	}
	return true, apperror.NewExecutionError(fmt.Sprintf("remediation failed for %s", f.RepoName))
}

func tryRemediateMatchedRecord(rec model.ScanRecord) (bool, error) {
	if !isSameOwnerHTTPSRepo(rec.AbsolutePath) {
		return false, nil
	}
	f := PullFailureSummary{
		RepoName: rec.RepoName,
		RepoPath: rec.AbsolutePath,
	}
	if remediateSingleRepo(f) {
		return true, nil
	}
	return true, apperror.NewExecutionError(fmt.Sprintf("remediation failed for %s", rec.RepoName))
}

// HasPendingPullFailures reports whether gitmap-pull.db has recorded pull failures.
func HasPendingPullFailures() bool {
	failures, err := loadRecentPullFailuresFromDB()
	return err == nil && len(failures) > 0
}

// ExtractPullFailures collects failures from slice of repo states.
func ExtractPullFailures(states []*PullRepoState) []PullFailureSummary {
	var failures []PullFailureSummary
	for _, s := range states {
		if s.ErrorMsg != "" || s.Step == PullStepTypeError {
			repoID, remoteURL, _ := LookupRepoRemoteAndID(s.RepoPath)
			failures = append(failures, PullFailureSummary{
				RepoID:    repoID,
				RepoName:  s.RepoName,
				RepoPath:  s.RepoPath,
				RemoteURL: remoteURL,
				ErrorText: s.ErrorMsg,
				ErrorType: string(s.Step),
			})
		}
	}

	return failures
}
