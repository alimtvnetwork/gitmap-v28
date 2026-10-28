package cmdpull

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// InactiveRepoDetail holds metadata for a repository skipped due to inactivity.
type InactiveRepoDetail struct {
	RepoName string `json:"repoName"`
	RepoPath string `json:"repoPath"`
	Reason   string `json:"reason"`
}

// EfficientPullPartition holds separated active scan records and inactive details.
type EfficientPullPartition struct {
	ActiveRecords []model.ScanRecord
	InactiveRepos []InactiveRepoDetail
}

// EfficientPullOptions contains runtime flags for efficient pull execution.
type EfficientPullOptions struct {
	IsTableMode  bool
	IsJSON       bool
	InvokedAlias string
	IsShortForm  bool
	UseSSH       bool
	UseHTTPS     bool
	TargetSSH    string
	WindowHours  int
	Args         []string
}

// PullRepoSummaryItem holds the outcome of pulling an individual repository.
type PullRepoSummaryItem struct {
	RepoName        string `json:"repoName"`
	Status          string `json:"status"`
	Changes         string `json:"changes"`
	ErrorDetails    string `json:"errorDetails,omitempty"`
	RemediationHint string `json:"remediationHint,omitempty"`
}

// PullEfficientSummary represents the serialized JSON summary of an efficient pull operation.
type PullEfficientSummary struct {
	Node          string                `json:"node,omitempty"`
	Total         int                   `json:"total"`
	ActiveCount   int                   `json:"activeCount"`
	InactiveCount int                   `json:"inactiveCount"`
	States        []PullRepoSummaryItem `json:"states"`
	Inactive      []InactiveRepoDetail  `json:"inactive"`
	DurationMs    int64                 `json:"durationMs"`
}

func resolveEfficientFullCmdName(isTable bool) string {
	if isTable {
		return "pull all-efficient-table"
	}
	return "pull all-efficient"
}

func handleEmptyRecords(isJSON bool) error {
	if isJSON {
		return renderJSONInactiveResults(0, nil)
	}
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %snothing to pull: no tracked repositories found in database.%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
	return nil
}

func printEfficientPartitionNotice(total, activeCount, inactiveCount int, windowHours ...int) {
	effectiveWindow := constants.DefaultPullActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s resolved %d repo(s) (%d active, %d inactive skipped in %dh window)\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		total, activeCount, inactiveCount, effectiveWindow)
}

func startActiveBatchHeartbeat(bar *PullProgressBar, isJSON bool) func() {
	return StartPullHeartbeat(30*time.Second, 5*time.Second, func() string {
		return fmt.Sprintf("%d/%d completed • %d active workers • %d failures so far...",
			bar.Completed(), bar.Total(), len(bar.ActiveWorkers()), bar.Failed())
	}, isJSON)
}

func launchThrottledActiveIgnoreScan(records []model.ScanRecord) *IgnoreScanHandle {
	ttl := resolvePullIgnoreTTL()
	cold, err := store.FilterReposNeedingIgnoreCheck(records, ttl)
	if err != nil || len(cold) == 0 {
		return nil
	}
	return StartThrottledAsyncIgnoreScan(cold, ttl)
}

func collectAndRemediateEfficientIgnore(handle *IgnoreScanHandle, args []string, isJSON bool) {
	if handle == nil {
		return
	}
	issues := handle.Collect()
	if len(issues) == 0 || isJSON {
		return
	}
	deduped := DeduplicateIgnoreIssues(issues)
	if len(deduped) == 0 {
		return
	}
	if hasEfficientAutoYes(args) {
		handleIgnoreRemediation(deduped, true, isJSON)
		return
	}
	markIssuesSkipped(deduped)
	printNonInteractiveIgnoreNotice(len(deduped))
}

func hasEfficientAutoYes(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "-y" || low == "--yes" || low == "--auto-fix" {
			return true
		}
	}
	return false
}

func handleEfficientBatchFinish(total int, bar *PullProgressBar, part EfficientPullPartition, all []model.ScanRecord, opts EfficientPullOptions, dur time.Duration) error {
	sortedStates := sortStatesAlphabetically(bar.States())
	if opts.IsJSON {
		return renderJSONEfficientResults(total, sortedStates, part.InactiveRepos, dur)
	}
	renderEfficientResults(sortedStates, part.InactiveRepos, opts.IsTableMode, all, opts.WindowHours)
	recordEfficientPullTelemetry(total, sortedStates, part.InactiveRepos, dur, opts.IsTableMode)

	unresolved := ExtractUnresolvedRepos(sortedStates)
	if len(unresolved) > 0 {
		_ = PromptPostPullRemediation(unresolved)
	}

	return nil
}

func renderEfficientResults(states []*PullRepoState, inactive []InactiveRepoDetail, isTable bool, allRecords []model.ScanRecord, windowHours ...int) {
	effectiveWindow := 0
	if len(windowHours) > 0 {
		effectiveWindow = windowHours[0]
	}
	if isTable {
		renderPullBatchResults(states)
	} else {
		renderConciseActiveResults(states, allRecords)
	}
	renderEfficientPullSummary(states, inactive, effectiveWindow)
}

func renderEfficientPullSummary(states []*PullRepoState, inactive []InactiveRepoDetail, windowHours ...int) {
	effectiveWindow := 0
	if len(windowHours) > 0 {
		effectiveWindow = windowHours[0]
	}
	fmt.Printf("\n  %s✓%s %sPull efficient complete:%s %d active pulled, %d inactive skipped\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		len(states), len(inactive))
	if len(inactive) > 0 {
		printInactiveSkipSummary(inactive, effectiveWindow)
	}
}

func printInactiveSkipSummary(inactive []InactiveRepoDetail, windowHours ...int) {
	effectiveWindow := constants.DefaultPullActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	arrow := resolveSubArrow()
	names := make([]string, 0, len(inactive))
	for _, in := range inactive {
		names = append(names, in.RepoName)
	}
	fmt.Printf("    %s%s%s %sskipped inactive repos (0 changes or checked recently in %dh window):%s\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, effectiveWindow, constants.ColorReset)
	termWidth := detectTerminalWidth()
	wrapped := FormatWrappedInactiveList(names, "      ", termWidth-2)
	fmt.Printf("%s%s%s\n", constants.ColorDim, wrapped, constants.ColorReset)
	fmt.Printf("    %sTo pull all repositories including inactive ones, run:%s %sgitmap pull all%s\n\n",
		constants.ColorDim, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
}

func handleAllReposInactive(total int, inactive []InactiveRepoDetail, isTable bool, isJSON bool, windowHours ...int) error {
	if isJSON {
		return renderJSONInactiveResults(total, inactive)
	}
	effectiveWindow := constants.DefaultPullActivityWindowHours
	if len(windowHours) > 0 && windowHours[0] > 0 {
		effectiveWindow = windowHours[0]
	}
	fmt.Printf("\n  %sℹ%s All %d repository(ies) are currently up-to-date (0 changes or checked recently in %dh window).\n",
		constants.ColorCyan, constants.ColorReset, total, effectiveWindow)
	printInactiveSkipSummary(inactive, effectiveWindow)
	recordEfficientPullTelemetry(total, nil, inactive, 0, isTable)
	return nil
}

func recordEfficientPullTelemetry(total int, states []*PullRepoState, inactive []InactiveRepoDetail, dur time.Duration, isTable bool) {
	cwd, _ := os.Getwd()
	cmdType := resolveEfficientFullCmdName(isTable)
	successCount, failedCount := 0, 0
	for _, s := range states {
		if s.ErrorMsg == "" && s.Step != PullStepTypeError {
			successCount++
		} else {
			failedCount++
		}
	}
	telemetry := PullSessionTelemetry{
		CommandType:   cmdType,
		WorkingDir:    cwd,
		TotalRepos:    total,
		PulledRepos:   len(states),
		SkippedRepos:  len(inactive),
		SuccessCount:  successCount,
		FailedCount:   failedCount,
		IsEfficient:   true,
		Duration:      dur,
		GitMapVersion: constants.Version,
	}
	_ = recordTelemetryToDB(telemetry, buildCombinedRepoRuns(states, inactive), states)
}

func buildCombinedRepoRuns(states []*PullRepoState, inactive []InactiveRepoDetail) []store.PullRepoRunRecord {
	runs := buildRepoRunRecords(states)
	for _, in := range inactive {
		runs = append(runs, store.PullRepoRunRecord{
			RepoPath:   in.RepoPath,
			RepoName:   in.RepoName,
			PullStatus: "skipped-inactive",
			IsActive:   false,
			HasChanges: false,
			DurationMs: 0,
			Notes:      in.Reason,
		})
	}
	return runs
}

func recordTelemetryToDB(telemetry PullSessionTelemetry, records []store.PullRepoRunRecord, states []*PullRepoState) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	runID, err := db.InsertPullRun(telemetry.ConvertToStoreRunRecord())
	if err != nil {
		return err
	}
	for _, s := range states {
		recordSinglePullErrorToDB(db, s)
	}
	return db.InsertPullRepoRuns(runID, records)
}

func recordSinglePullErrorToDB(db *store.PullSplitDB, s *PullRepoState) {
	if s == nil {
		return
	}
	if s.ErrorMsg != "" || s.Step == PullStepTypeError {
		rec := store.PullErrorRecord{
			RepoSlug:       s.RepoName,
			RepoPath:       s.RepoPath,
			NodeID:         "local-01",
			NodeVersion:    constants.Version,
			ErrorType:      ResolvePullErrorDetails(s),
			ErrorText:      s.ErrorMsg,
			RemediationCmd: ResolvePullRemediationHint(s),
			CreatedAt:      time.Now().UTC(),
		}
		if insErr := db.InsertPullError(rec); insErr != nil {
			fmt.Fprintf(os.Stderr, "  %s⚠ Error inserting pull error for %s: %v%s\n",
				constants.ColorYellow, s.RepoName, insErr, constants.ColorReset)
		}
		return
	}
	if clearErr := db.ClearPullErrorsForRepo(s.RepoName, s.RepoPath); clearErr != nil {
		fmt.Fprintf(os.Stderr, "  %s⚠ Notice: error clearing pull status for %s: %v%s\n",
			constants.ColorYellow, s.RepoName, clearErr, constants.ColorReset)
	}
}
