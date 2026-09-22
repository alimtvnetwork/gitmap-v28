package cmdpipeline

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// PipelineDetailsFlags holds options for gitmap pipeline details.
type PipelineDetailsFlags struct {
	IsJSON       bool
	IsDetailed   bool
	HasForce     bool
	HasIndex     bool
	Index        int
	TargetRunId  uint64
	CommitTarget string
}

// PipelineDetailsPayload encapsulates the complete job details and error telemetry for a workflow run.
type PipelineDetailsPayload struct {
	Repo            string                  `json:"repo"`
	WorkflowName    string                  `json:"workflowName"`
	RunId           uint64                  `json:"runId"`
	Branch          string                  `json:"branch"`
	Sha             string                  `json:"sha"`
	Status          string                  `json:"status"`
	Conclusion      string                  `json:"conclusion"`
	DurationSeconds int                     `json:"durationSeconds"`
	IsFromCache     bool                    `json:"isFromCache,omitempty"`
	CacheSource     string                  `json:"cacheSource,omitempty"`
	Jobs            []PipelineJobDetailItem `json:"jobs"`
	FailedCount     int                     `json:"failedCount"`
	PassedCount     int                     `json:"passedCount"`
	RunningCount    int                     `json:"runningCount"`
	TotalCount      int                     `json:"totalCount"`
	Failures        []SectionFailure        `json:"failures,omitempty"`
	ErrorLogs       string                  `json:"errorLogs,omitempty"`
}

// PipelineJobDetailItem represents a single runner target job in the details table.
type PipelineJobDetailItem struct {
	RunnerTarget    string `json:"runnerTarget"`
	Step            string `json:"step"`
	ElapsedTime     string `json:"elapsedTime"`
	DurationSeconds int    `json:"durationSeconds"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	StatusDisplay   string `json:"statusDisplay"`
	FailureSummary  string `json:"failureSummary,omitempty"`
}

func parsePipelineDetailsFlags(args []string) PipelineDetailsFlags {
	var flags PipelineDetailsFlags
	flags.IsJSON = hasArgFlag(args, "--json")
	flags.IsDetailed = hasDetailedArg(args)
	flags.HasForce = hasForceArg(args)
	flags.Index, flags.HasIndex = scanNegativeIndex(args)
	flags.TargetRunId = scanTargetRunId(args)
	flags.CommitTarget = scanCommitTarget(args)

	return flags
}

func scanTargetRunId(args []string) uint64 {
	for _, a := range args {
		if val, err := strconv.ParseUint(a, 10, 64); err == nil && val > 100000 {
			return val
		}
	}

	return 0
}

func scanCommitTarget(args []string) string {
	for _, a := range args {
		if isCommitShaCandidate(a) {
			return a
		}
	}

	return ""
}

func isCommitShaCandidate(a string) bool {
	if strings.HasPrefix(a, "-") {
		return false
	}
	clean := strings.TrimSpace(a)
	if len(clean) >= 7 && len(clean) <= 40 && isAllHex(clean) {
		return true
	}

	return false
}

func isAllHex(s string) bool {
	for _, r := range strings.ToLower(s) {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}

	return true
}

// HandlePipelineDetails renders the detailed runner targets table and failure summary.
func HandlePipelineDetails(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineDetailsHelp()
		return nil
	}

	return executePipelineDetails(parsePipelineDetailsFlags(args))
}

func executePipelineDetails(flags PipelineDetailsFlags) error {
	repo := resolveCurrentRepoSlug()
	targetRun, isFromCache := resolveDetailsTargetRun(repo, flags)
	if targetRun.DatabaseId == 0 {
		fmt.Printf("No pipeline runs found for %s.\n", repo)
		return nil
	}

	jobs := queryRunJobs(repo, targetRun.DatabaseId)
	payload := buildPipelineDetailsPayload(repo, targetRun, jobs, isFromCache)

	return dispatchDetailsPresentation(payload, flags)
}

func dispatchDetailsPresentation(payload PipelineDetailsPayload, flags PipelineDetailsFlags) error {
	if flags.IsJSON {
		return printJSON(payload)
	}

	renderPipelineDetailsTerminal(payload)

	return nil
}

func resolveDetailsTargetRun(repo string, flags PipelineDetailsFlags) (ghRunItem, bool) {
	decision := EvaluatePipelineErrorsCache(repo, PipelineErrorFlags{
		HasForce: flags.HasForce,
		HasIndex: flags.HasIndex,
		Index:    flags.Index,
	})
	if decision.IsFromCache && len(decision.CachedRuns) > 0 {
		return selectRunByIndex(decision.CachedRuns, flags), true
	}

	runs := queryWorkflowRuns(repo)
	RecordFetchedRunsToSplitDb(repo, runs)

	return selectRunByIndex(runs, flags), false
}

func selectRunByIndex(runs []ghRunItem, flags PipelineDetailsFlags) ghRunItem {
	if len(runs) == 0 {
		return ghRunItem{}
	}
	if flags.HasIndex && flags.Index < 0 {
		return resolveNegativeOffsetRun(runs, flags.Index)
	}

	return findPrimaryTargetRun(runs)
}

func resolveNegativeOffsetRun(runs []ghRunItem, index int) ghRunItem {
	offset := NormalizeNegativeIndex(index)
	if offset < len(runs) {
		return runs[offset]
	}

	return findPrimaryTargetRun(runs)
}

func initBaseDetailsPayload(repo string, run ghRunItem, isFromCache bool) PipelineDetailsPayload {
	return PipelineDetailsPayload{
		Repo:            repo,
		WorkflowName:    run.Name,
		RunId:           run.DatabaseId,
		Branch:          run.HeadBranch,
		Sha:             run.HeadSha,
		Status:          run.Status,
		Conclusion:      run.Conclusion,
		DurationSeconds: calculateRunDuration(run.CreatedAt, run.UpdatedAt),
		IsFromCache:     isFromCache,
		Jobs:            make([]PipelineJobDetailItem, 0),
		Failures:        make([]SectionFailure, 0),
	}
}

func buildPipelineDetailsPayload(repo string, run ghRunItem, jobs []ghJobItem, isFromCache bool) PipelineDetailsPayload {
	payload := initBaseDetailsPayload(repo, run, isFromCache)
	populateJobsAndSummary(&payload, repo, jobs)
	attachFailuresToDetailsPayload(&payload, repo)

	return payload
}

func populateJobsAndSummary(p *PipelineDetailsPayload, repo string, jobs []ghJobItem) {
	p.TotalCount = len(jobs)
	for _, j := range jobs {
		item := buildJobDetailItem(repo, j)
		p.Jobs = append(p.Jobs, item)
		tallyJobOutcome(p, item)
	}
}

func tallyJobOutcome(p *PipelineDetailsPayload, item PipelineJobDetailItem) {
	if item.Conclusion == "success" {
		p.PassedCount++
		return
	}
	if item.Conclusion == "failure" {
		p.FailedCount++
		return
	}
	if item.Status == "in_progress" || item.Status == "queued" {
		p.RunningCount++
	}
}

func buildJobDetailItem(repo string, j ghJobItem) PipelineJobDetailItem {
	durText, durSec := formatJobElapsedTime(j)
	statusDisplay, summary := formatJobStatusDisplay(j)

	return PipelineJobDetailItem{
		RunnerTarget:    formatRunnerTargetName(j.Name),
		Step:            resolveJobActiveOrFailingStep(j),
		ElapsedTime:     durText,
		DurationSeconds: durSec,
		Status:          j.Status,
		Conclusion:      j.Conclusion,
		StatusDisplay:   statusDisplay,
		FailureSummary:  summary,
	}
}

func formatRunnerTargetName(name string) string {
	cleaned := strings.TrimSpace(name)
	if len(cleaned) == 0 {
		return "Runner"
	}

	return cleaned
}

func resolveJobActiveOrFailingStep(j ghJobItem) string {
	if j.Conclusion == "failure" {
		return resolveFailedJobStep(j.Steps)
	}
	if j.Status == "in_progress" {
		return resolveActiveJobStep(j.Steps)
	}

	return findSubstantiveStepName(j.Steps)
}

func resolveFailedJobStep(steps []ghStepItem) string {
	failStep := findFailingStepName(steps)
	if len(failStep) > 0 {
		return failStep
	}

	return findSubstantiveStepName(steps)
}

func resolveActiveJobStep(steps []ghStepItem) string {
	activeStep := findActiveStepName(steps)
	if len(activeStep) > 0 {
		return activeStep
	}

	return findSubstantiveStepName(steps)
}

func findFailingStepName(steps []ghStepItem) string {
	for _, s := range steps {
		if s.Conclusion == "failure" {
			return s.Name
		}
	}

	return ""
}

func findActiveStepName(steps []ghStepItem) string {
	for _, s := range steps {
		if s.Status == "in_progress" {
			return s.Name
		}
	}

	return ""
}

func findSubstantiveStepName(steps []ghStepItem) string {
	for i := len(steps) - 1; i >= 0; i-- {
		s := steps[i]
		if !strings.HasPrefix(s.Name, "Post ") && !strings.HasPrefix(s.Name, "Complete ") {
			return s.Name
		}
	}
	if len(steps) > 0 {
		return steps[len(steps)-1].Name
	}

	return "Execute job"
}

func formatJobElapsedTime(j ghJobItem) (string, int) {
	durSec := computeJobDuration(j.StartedAt, j.CompletedAt, j.Status)
	if durSec <= 0 {
		return "-", 0
	}
	if durSec < 60 {
		return fmt.Sprintf("~%ds", durSec), durSec
	}

	m := durSec / 60
	return fmt.Sprintf("~%ds (%dm)", durSec, m), durSec
}

func computeJobDuration(startedStr, completedStr, status string) int {
	t1, err := time.Parse(time.RFC3339, startedStr)
	if err != nil {
		return 0
	}
	if status == "in_progress" {
		return int(time.Since(t1).Seconds())
	}

	return computeFinishedJobDuration(t1, completedStr)
}

func computeFinishedJobDuration(t1 time.Time, completedStr string) int {
	t2, err := time.Parse(time.RFC3339, completedStr)
	if err != nil || t2.Before(t1) {
		return 0
	}

	return int(t2.Sub(t1).Seconds())
}

func formatSuccessStatusDisplay(name string) string {
	note := extractPackagingNote(name)
	if len(note) > 0 {
		return fmt.Sprintf("%sPASSED ✔%s (%s)", constants.ColorGreen, constants.ColorReset, note)
	}

	return fmt.Sprintf("%sPASSED ✔%s", constants.ColorGreen, constants.ColorReset)
}

func formatFailureStatusDisplay(steps []ghStepItem) (string, string) {
	failStep := findFailingStepName(steps)
	if len(failStep) > 0 {
		return fmt.Sprintf("%sFAILED ✖%s (%s)", constants.ColorRed, constants.ColorReset, failStep), failStep
	}

	return fmt.Sprintf("%sFAILED ✖%s", constants.ColorRed, constants.ColorReset), ""
}

func formatJobStatusDisplay(j ghJobItem) (string, string) {
	if j.Conclusion == "success" {
		return formatSuccessStatusDisplay(j.Name), ""
	}
	if j.Conclusion == "failure" {
		return formatFailureStatusDisplay(j.Steps)
	}

	return formatActiveStatusDisplay(j.Status), ""
}

func formatActiveStatusDisplay(status string) string {
	if status == "in_progress" {
		return fmt.Sprintf("%sRUNNING ●%s", constants.ColorYellow, constants.ColorReset)
	}

	return fmt.Sprintf("%sQUEUED ○%s", constants.ColorDim, constants.ColorReset)
}

func extractPackagingNote(jobName string) string {
	lower := strings.ToLower(jobName)
	if strings.Contains(lower, "deb") || strings.Contains(lower, "ubuntu") {
		return "packaging AppImage/deb"
	}
	if strings.Contains(lower, "win") || strings.Contains(lower, "nsis") {
		return "packaging NSIS setup.exe"
	}

	return extractAppleOrArmPackagingNote(lower)
}

func extractAppleOrArmPackagingNote(lower string) string {
	if strings.Contains(lower, "aarch64") || strings.Contains(lower, "arm") {
		return "packaging app bundle"
	}
	if strings.Contains(lower, "mac") || strings.Contains(lower, "darwin") {
		return "packaging universal bundle"
	}

	return ""
}

func attachFailuresToDetailsPayload(p *PipelineDetailsPayload, repo string) {
	if p.FailedCount == 0 && p.Conclusion != "failure" {
		return
	}

	rawLogs := queryFailedRunLogs(repo, p.RunId)
	if len(rawLogs) == 0 {
		return
	}

	p.ErrorLogs = rawLogs
	p.Failures = buildSingleRunSectionFailures(repo, p)
}

func buildSingleRunSectionFailures(repo string, p *PipelineDetailsPayload) []SectionFailure {
	runItem := ghRunItem{
		DatabaseId: p.RunId, Name: p.WorkflowName, HeadBranch: p.Branch,
		HeadSha: p.Sha, Conclusion: p.Conclusion, Status: p.Status,
	}

	return extractAllSectionFailures([]FailedRunItem{fetchAndBuildFailedRunItem(repo, runItem)})
}

func renderPipelineDetailsTerminal(p PipelineDetailsPayload) {
	fmt.Printf("\n  %s● Pipeline Details: %s%s\n\n", constants.ColorCyan, p.Repo, constants.ColorReset)
	renderDetailsSummaryBlock(p)
	fmt.Println()
	renderDetailsTable(p.Jobs)
	fmt.Println()
	renderDetailsFailuresBlock(p)
}

func renderDetailsSummaryBlock(p PipelineDetailsPayload) {
	fmt.Printf("    %-16s %s #%d\n", "Workflow:", p.WorkflowName, p.RunId)
	fmt.Printf("    %-16s %s\n", "Branch:", p.Branch)
	fmt.Printf("    %-16s %s\n", "Last Commit:", gitutil.TruncSha(p.Sha))
	fmt.Printf("    %-16s %s\n", "Run Status:", formatDetailsRunStatus(p.Status, p.Conclusion))
	fmt.Printf("    %-16s %s\n", "Duration:", formatDurationSeconds(p.DurationSeconds))
	fmt.Printf("    %-16s %s\n", "Jobs Summary:", formatJobsCountsSummary(p))
	if p.IsFromCache {
		fmt.Printf("    %-16s %s⚡ Served from local %s DB cache (commit %s)%s\n",
			"Cache:", constants.ColorCyan, strings.ToUpper(p.CacheSource), gitutil.TruncSha(p.Sha), constants.ColorReset)
	}
}

func formatDetailsRunStatus(status, conclusion string) string {
	if conclusion == "success" {
		return constants.ColorGreen + "PASSING (100% green)" + constants.ColorReset
	}
	if conclusion == "failure" {
		return constants.ColorRed + "FAILING (errors detected)" + constants.ColorReset
	}
	if status == "in_progress" {
		return constants.ColorYellow + "RUNNING (in progress)" + constants.ColorReset
	}

	return status
}

func formatJobsCountsSummary(p PipelineDetailsPayload) string {
	return fmt.Sprintf("%d Passed | %d Failed | %d Running (%d total)",
		p.PassedCount, p.FailedCount, p.RunningCount, p.TotalCount)
}

func renderDetailsTable(jobs []PipelineJobDetailItem) {
	printDetailsTableHeader()
	for _, j := range jobs {
		renderDetailsTableRow(j)
	}
}

func printDetailsTableHeader() {
	header := fmt.Sprintf("    %-24s  %-24s  %-16s  %s",
		"Runner Target", "Step", "Elapsed Time", "Status")
	divider := fmt.Sprintf("    %-24s  %-24s  %-16s  %s",
		strings.Repeat("─", 24), strings.Repeat("─", 24), strings.Repeat("─", 16), strings.Repeat("─", 28))
	fmt.Println(header)
	fmt.Println(divider)
}

func renderDetailsTableRow(j PipelineJobDetailItem) {
	target := truncateHistoryStr(j.RunnerTarget, 24)
	step := truncateHistoryStr(j.Step, 24)
	elapsed := truncateHistoryStr(j.ElapsedTime, 16)
	paddedTarget := padRightVisible(target, 24)
	paddedStep := padRightVisible(step, 24)
	paddedElapsed := padRightVisible(elapsed, 16)

	fmt.Printf("    %s  %s  %s  %s\n", paddedTarget, paddedStep, paddedElapsed, j.StatusDisplay)
}

func renderDetailsFailuresBlock(p PipelineDetailsPayload) {
	if len(p.Failures) == 0 {
		return
	}

	fmt.Printf("  %s● Failure Diagnostics Across Runner Targets [%d failed step(s)]:%s\n\n",
		constants.ColorRed, len(p.Failures), constants.ColorReset)
	for i, f := range p.Failures {
		renderSingleDetailsFailure(f, i+1, len(p.Failures))
	}

	fmt.Printf("  💡 Suggested action: %sgitmap aef%s (or: %sgitmap pe%s)\n\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
}

func renderSingleDetailsFailure(sec SectionFailure, idx, total int) {
	fmt.Printf("    %s[%d/%d] %s ➔ Job: %s | Step: %s%s\n",
		constants.ColorCyan, idx, total, sec.WorkflowName, sec.JobName, sec.StepName, constants.ColorReset)
	if len(sec.FailureSummary) > 0 {
		fmt.Printf("      Error:   %s%s%s\n", constants.ColorRed, sec.FailureSummary, constants.ColorReset)
	}
	renderSectionWarnings(sec.Warnings)
	renderSectionErrorLinesDedup(sec.ErrorLines, sec.FailureSummary)
	if len(sec.SavedLogFile) > 0 {
		fmt.Printf("      Log:     %s\n", filepath.ToSlash(sec.SavedLogFile))
	}
}

func printPipelineDetailsHelp() {
	fmt.Println("Usage: gitmap pipeline details [flags]")
	fmt.Println("       gitmap pd [flags] (shortcut for pipeline details)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --json                  Output data in structured JSON format")
	fmt.Println("  --force, --no-cache     Bypass local SQLite DB cache and pull fresh from GitHub")
	fmt.Println("  -v, --detailed          Show full uncompacted step error traces")
	fmt.Println("  -1, -2, -3              Inspect details for past commit by negative offset")
}
