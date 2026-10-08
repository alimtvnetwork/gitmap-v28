package cmdpull

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func runPullBatchExecution(records []model.ScanRecord, opts pullOptions) (*PullProgressBar, []*PullRepoState, time.Duration) {
	bar := NewPullProgressBar(len(records), opts.isJSON, opts.stopOnFail)
	if !opts.isJSON {
		bar.Start()
	}
	startTime := time.Now()
	executePull(records, bar, opts)
	if !opts.isJSON {
		bar.Stop()
	}
	dur := time.Since(startTime)
	dedupedStates := DeduplicateRepoStates(bar.States())

	return bar, sortStatesAlphabetically(dedupedStates), dur
}

func renderPullBatchOutput(records []model.ScanRecord, sortedStates []*PullRepoState, dur time.Duration, opts pullOptions) {
	if isConcisePullOutput(opts.all, opts.showStatus) {
		renderConciseActiveResults(sortedStates, records)
		activeCount, upToDateCount := countActiveStates(sortedStates)
		printPullAllFastSummary(len(sortedStates), activeCount, upToDateCount, dur)
	} else {
		renderPullBatchResults(sortedStates)
	}
	handlePullRemediationForRecords(records, opts)
}

func isConcisePullOutput(isAll, isShowStatus bool) bool {
	if isShowStatus {
		return false
	}

	return isAll
}

func countActiveStates(states []*PullRepoState) (int, int) {
	activeCount := 0
	upToDateCount := 0
	for _, s := range states {
		if ResolveRepoStatusLabel(s.Changes) == "up-to-date" {
			upToDateCount++
		} else {
			activeCount++
		}
	}

	return activeCount, upToDateCount
}

func printPullAllFastSummary(pulledCount, activeCount, upToDateCount int, dur time.Duration) {
	if activeCount > 0 {
		fmt.Printf("\n  %s✓%s %sPull all complete:%s %d pulled (%d active, %d up-to-date) (%s)\n\n",
			constants.ColorGreen, constants.ColorReset,
			constants.ColorBold, constants.ColorReset,
			pulledCount, activeCount, upToDateCount, FormatPullDuration(dur))

		return
	}
	fmt.Printf("\n  %s✓%s %sPull all complete:%s %d pulled (%s)\n\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		pulledCount, FormatPullDuration(dur))
}

func syncPullBatchTelemetry(records []model.ScanRecord, states []*PullRepoState, dur time.Duration, opts pullOptions) {
	telemetry := buildPullSessionTelemetry(records, states, dur, opts.all)
	_ = RecordPullBatchSession(telemetry, states)
}

func buildPullSessionTelemetry(records []model.ScanRecord, states []*PullRepoState, dur time.Duration, isAll bool) PullSessionTelemetry {
	cwd, _ := os.Getwd()
	succ, fail := countBatchStateOutcomes(states)

	return PullSessionTelemetry{
		CommandType: resolveBatchCommandType(isAll), WorkingDir: cwd,
		TotalRepos: len(records), PulledRepos: len(states),
		SkippedRepos: len(records) - len(states), SuccessCount: succ, FailedCount: fail,
		IsEfficient: false, Duration: dur, GitMapVersion: constants.Version,
	}
}

func resolveBatchCommandType(isAll bool) string {
	if isAll {
		return "pull-all"
	}

	return "pull"
}

func countBatchStateOutcomes(states []*PullRepoState) (int, int) {
	successCount := 0
	failedCount := 0
	for _, s := range states {
		if s.ErrorMsg == "" && s.Step != PullStepTypeError {
			successCount++
		} else {
			failedCount++
		}
	}

	return successCount, failedCount
}

func maybeApplyTransportToRecords(records []model.ScanRecord, useSSH, useHTTPS bool) {
	hasTransport := useSSH || useHTTPS
	if !hasTransport {
		return
	}
	for _, rec := range records {
		applyTransportIfPathExists(rec.AbsolutePath, useSSH, useHTTPS)
	}
}

func applyTransportIfPathExists(path string, useSSH, useHTTPS bool) {
	hasPath := len(path) > 0
	if !hasPath {
		return
	}

	ApplyTransportFlag(path, useSSH, useHTTPS)
}

func sortStatesAlphabetically(states []*PullRepoState) []*PullRepoState {
	sorted := make([]*PullRepoState, len(states))
	copy(sorted, states)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})

	return sorted
}

func renderPullBatchResults(states []*PullRepoState) {
	var tableRows []model.PullTableRow
	for _, state := range states {
		tableRows = append(tableRows, buildPullTableRowFromState(state))
	}
	RenderPullBatchTable(tableRows)
}

func buildPullTableRowFromState(state *PullRepoState) model.PullTableRow {
	row := initBasePullTableRow(state)
	row.LatestBranch = gitutil.GetLatestRemoteBranch(state.RepoPath)
	row.PRStatus = gitutil.DetectPRStatusFast(state.RepoPath)
	row.PullStatus = resolveStateStatus(state)
	row.Release = resolveRepoRelease(state.RepoPath, row.LatestBranch)
	row.LastSHA = resolveRowSHA(state)

	return row
}

func resolveRowSHA(state *PullRepoState) string {
	sha := state.NewSHA
	if len(sha) == 0 {
		sha = state.OldSHA
	}
	if len(sha) == 0 && len(state.RepoPath) > 0 {
		sha = gitutil.GetLastCommitSHA(state.RepoPath)
	}

	return ShortenSHA(sha)
}

func resolveRepoRelease(repoPath, latestBranch string) string {
	hasPath := len(repoPath) > 0
	if !hasPath {
		return "-"
	}

	rel := queryReleaseIdentifier(repoPath, latestBranch)
	if len(rel) > 0 {
		return rel
	}

	return "-"
}

func queryReleaseIdentifier(repoPath, latestBranch string) string {
	tag := gitutil.GetLatestTag(repoPath)
	if len(tag) > 0 {
		return tag
	}

	if rel := extractReleaseFromBranch(latestBranch); len(rel) > 0 {
		return rel
	}

	return readRepoManifestVersion(repoPath)
}

func extractReleaseFromBranch(branch string) string {
	lower := strings.ToLower(branch)
	if strings.HasPrefix(lower, "release/") {
		return strings.TrimPrefix(branch, "release/")
	}
	if strings.HasPrefix(lower, "v") && isSemverLike(branch) {
		return branch
	}

	return ""
}

func isSemverLike(s string) bool {
	raw := strings.TrimPrefix(s, "v")
	parts := strings.Split(raw, ".")
	hasEnoughParts := len(parts) >= 2
	if !hasEnoughParts {
		return false
	}

	_, err := strconv.Atoi(parts[0])

	return err == nil
}

func readRepoManifestVersion(repoPath string) string {
	v := readJsonVersion(filepath.Join(repoPath, "version.json"))
	if len(v) > 0 {
		return ensureVPrefix(v)
	}
	pkgV := readJsonVersion(filepath.Join(repoPath, "package.json"))
	if len(pkgV) > 0 {
		return ensureVPrefix(pkgV)
	}

	return ""
}

func readJsonVersion(filePath string) string {
	data, err := os.ReadFile(filePath)
	hasData := err == nil && len(data) > 0
	if !hasData {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}

	return extractVersionField(m)
}

func extractVersionField(m map[string]interface{}) string {
	if v, ok := m["Version"].(string); ok && len(v) > 0 {
		return v
	}
	if v, ok := m["version"].(string); ok && len(v) > 0 {
		return v
	}

	return ""
}

func ensureVPrefix(v string) string {
	trimmed := strings.TrimSpace(v)
	hasV := strings.HasPrefix(trimmed, "v")
	if hasV {
		return trimmed
	}

	return "v" + trimmed
}

func initBasePullTableRow(state *PullRepoState) model.PullTableRow {
	return model.PullTableRow{
		RepoName:    state.RepoName,
		Branch:      state.Branch,
		LastSHA:     ShortenSHA(state.NewSHA),
		CommitRange: state.CommitRange,
		Changes:     state.Changes,
		Duration:    FormatPullDuration(state.Duration),
		IsDirty:     state.IsDirty,
		Reason:      state.ErrorMsg,
	}
}

func resolveStateStatus(state *PullRepoState) string {
	if state.IsDirty {
		return "DIRTY"
	}

	return string(state.Step)
}
