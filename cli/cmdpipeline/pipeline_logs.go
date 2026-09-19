package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/atotto/clipboard"
)

func handlePipelineErrorLogs(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineErrorLogsHelp()

		return nil
	}
	if hasArgFlag(args, "last-failed-logs") {
		return HandlePipelineLastFailedLogs(args)
	}

	return handlePipelineHistoryOrExecute(args)
}

func handlePipelineHistoryOrExecute(args []string) error {
	if handled, err := HandlePipelineHistoryErrors(args); handled {
		return err
	}

	return executePipelineErrorLogs(args)
}

func executePipelineErrorLogs(args []string) error {
	flags := ParsePipelineErrorFlags(args)
	repo := resolveCurrentRepoSlug()
	if flags.HasTimeline {
		return executeTimelineErrorLogs(repo, flags, args)
	}

	return processAndRenderErrorLogs(repo, flags)
}

func executeTimelineErrorLogs(repo string, flags PipelineErrorFlags, args []string) error {
	WaitForRunnerETAIfActive()

	return runPipelineErrorLogsDynamicTimeline(ErrorLogsTimelineParams{
		Repo: repo, IsJSON: flags.IsJSON, WantFix: flags.HasFix,
		WantCheck: flags.HasCheck, IsDetailed: flags.IsDetailed,
		FilePath: flags.FilePath, TempFileName: flags.TempFileName, Args: args,
	})
}

func processAndRenderErrorLogs(repo string, flags PipelineErrorFlags) error {
	printReadingProgress(flags)
	runs := queryWorkflowRuns(repo)
	payload := buildErrorLogsPayload(repo, runs)
	applyPayloadOptions(&payload, runs, flags)

	return writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:              payload,
		IsJSON:               flags.IsJSON,
		HasSuppressOutputLog: flags.HasSuppressOutputLog,
		FilePath:             flags.FilePath,
		TempFile:             flags.TempFileName,
	})
}

func printReadingProgress(flags PipelineErrorFlags) {
	if flags.IsJSON || flags.HasSuppressOutputLog {
		return
	}

	fmt.Println("Reading pipeline logs...")
}

func applyPayloadOptions(p *PipelineErrorLogsPayload, runs []ghRunItem, flags PipelineErrorFlags) {
	if len(runs) > 0 {
		p.RerunEtaSeconds = calculateAverageDuration(runs, p.WorkflowName)
	}
	if !flags.IsDetailed {
		compactErrorPayload(p)
	}
	if flags.HasFix || flags.HasCheck {
		p.CICDChecks = runInternalCICDChecks(flags.HasFix)
	}
}

func buildErrorLogsPayload(repo string, runs []ghRunItem) PipelineErrorLogsPayload {
	payload := initBaseErrorLogsPayload(repo)
	if len(runs) == 0 {
		return handleEmptyRunsPayload(repo, runs, payload)
	}

	populateRunsIntoPayload(repo, runs, &payload)
	enrichErrorLogsMetadata(&payload, repo, runs)

	return payload
}

func handleEmptyRunsPayload(repo string, runs []ghRunItem, p PipelineErrorLogsPayload) PipelineErrorLogsPayload {
	if ApplyPreviousDbFallbackToPayload(&p, repo) {
		enrichErrorLogsMetadata(&p, repo, runs)

		return p
	}
	enrichErrorLogsMetadata(&p, repo, runs)

	return buildLocalOrEmptyErrorPayload(p)
}

func populateRunsIntoPayload(repo string, runs []ghRunItem, p *PipelineErrorLogsPayload) {
	initLatestRunMeta(p, runs[0])
	checkAndApplyRunningState(p, runs)
	failedRuns := resolveFailedRunsForPayload(repo, runs)
	if len(failedRuns) > 0 {
		populateFailedRunsPayload(repo, failedRuns, p)

		return
	}
	if shouldSkipFallback(p, runs) {
		return
	}
	_ = ApplyPreviousRunFallbackToPayload(p, repo, runs)
}

func shouldSkipFallback(p *PipelineErrorLogsPayload, runs []ghRunItem) bool {
	if hasFailingRuns(runs) {
		return false
	}

	return p.Conclusion == "success" || p.IsRunning
}

func hasFailingRuns(runs []ghRunItem) bool {
	for _, r := range runs {
		if isFailingConclusion(r.Conclusion) {
			return true
		}
	}

	return false
}

func isFailingConclusion(conclusion string) bool {
	switch strings.ToLower(strings.TrimSpace(conclusion)) {
	case "failure", "timed_out", "cancelled", "startup_failure":
		return true
	default:
		return false
	}
}

func enrichErrorLogsMetadata(p *PipelineErrorLogsPayload, repo string, runs []ghRunItem) {
	p.RepoUrl = resolveRepoWebURL(repo)
	p.LastReleaseVersion = queryLatestTagRelease(repo)
	p.OpenPRsCount = queryPendingPRs(repo)
	p.LatestBranch = resolveLatestBranchName(p, runs)
	p.LastHash = resolveLatestCommitHash(p, runs)
}

func resolveRepoWebURL(repo string) string {
	remoteURL, err := gitutil.RemoteURL(".")
	if err == nil && len(remoteURL) > 0 {
		return formatWebURL(remoteURL)
	}

	if len(repo) > 0 {
		return "https://github.com/" + repo
	}

	return ""
}

func formatWebURL(raw string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(raw), ".git")
	if strings.HasPrefix(clean, "git@github.com:") {
		return "https://github.com/" + strings.TrimPrefix(clean, "git@github.com:")
	}

	return clean
}

func resolveLatestBranchName(p *PipelineErrorLogsPayload, runs []ghRunItem) string {
	if len(runs) > 0 && len(runs[0].HeadBranch) > 0 {
		return runs[0].HeadBranch
	}
	if len(p.Branch) > 0 {
		return p.Branch
	}

	return resolveActiveOrMainBranch()
}

func resolveActiveOrMainBranch() string {
	active := gitutil.GetActiveBranch(".")
	if len(active) > 0 && active != "-" {
		return active
	}

	return "main"
}

func resolveLatestCommitHash(p *PipelineErrorLogsPayload, runs []ghRunItem) string {
	if len(runs) > 0 && len(runs[0].HeadSha) > 0 {
		return gitutil.TruncSha(runs[0].HeadSha)
	}
	if len(p.Sha) > 0 {
		return gitutil.TruncSha(p.Sha)
	}

	return resolveLocalCommitSHA()
}

func resolveLocalCommitSHA() string {
	sha := gitutil.GetLastCommitSHA(".")
	if len(sha) > 0 && sha != "-" {
		return sha
	}

	return ""
}

func initBaseErrorLogsPayload(repo string) PipelineErrorLogsPayload {
	return PipelineErrorLogsPayload{
		Repo:            repo,
		DbPath:          pipelinedb.PipelineDbPath(repo),
		SavedReportFile: resolvePipelineErrorReportPath(),
	}
}

func checkAndApplyRunningState(p *PipelineErrorLogsPayload, runs []ghRunItem) {
	if len(runs) == 0 {
		return
	}

	latest := runs[0]
	if latest.Status == "in_progress" || latest.Status == "queued" {
		setPayloadRunningState(p, latest, calculateETA(runs))
	}
}

func initLatestRunMeta(payload *PipelineErrorLogsPayload, latest ghRunItem) {
	payload.WorkflowName = latest.Name
	payload.RunId = latest.DatabaseId
	payload.Status = latest.Status
	payload.Conclusion = latest.Conclusion
	payload.Url = latest.Url
	payload.Branch = latest.HeadBranch
	payload.Sha = latest.HeadSha
	payload.CreatedAt = latest.CreatedAt
	payload.UpdatedAt = latest.UpdatedAt
	payload.DurationSeconds = calculateRunDuration(latest.CreatedAt, latest.UpdatedAt)
	payload.SavedLogFile = getCachedPipelineLogPath(latest.DatabaseId)
}

func setPayloadRunningState(p *PipelineErrorLogsPayload, r ghRunItem, eta int) {
	p.IsRunning = true
	p.ActiveRunName = r.Name
	p.ActiveRunId = r.DatabaseId
	p.ActiveRunUrl = r.Url
	p.EtaSeconds = eta
	p.ErrorLogs = fmt.Sprintf("Pipeline [%s #%d] is currently running. Estimated completion in %d seconds.",
		r.Name, r.DatabaseId, eta)
}

func resolveFailedRunsForPayload(repo string, runs []ghRunItem) []ghRunItem {
	return collectFailedRuns(runs)
}

func populateFailedRunsPayload(repo string, failedRuns []ghRunItem, p *PipelineErrorLogsPayload) {
	initFailedRunTopLevel(p, failedRuns[0])
	p.FailedRuns = fetchAllFailedRunsParallel(repo, failedRuns)
	p.SectionFailures = extractAllSectionFailures(p.FailedRuns)
	p.CombinedErrors = formatCombinedSectionFailures(p.SectionFailures)
	p.ErrorLogs = formatAggregatedErrorLogs(p.FailedRuns)
}

func fetchAllFailedRunsParallel(repo string, failedRuns []ghRunItem) []FailedRunItem {
	total := len(failedRuns)
	if total == 0 {
		return nil
	}
	if total == 1 {
		return []FailedRunItem{fetchAndBuildFailedRunItem(repo, failedRuns[0])}
	}

	return executeParallelFetchWorkers(repo, failedRuns)
}

func executeParallelFetchWorkers(repo string, failedRuns []ghRunItem) []FailedRunItem {
	results := make([]FailedRunItem, len(failedRuns))
	sem := make(chan struct{}, resolveFetchConcurrency(len(failedRuns)))
	var wg sync.WaitGroup

	for i, fr := range failedRuns {
		wg.Add(1)
		go dispatchFetchRunWorker(&wg, sem, results, repo, fr, i)
	}
	wg.Wait()

	return results
}

func dispatchFetchRunWorker(wg *sync.WaitGroup, sem chan struct{}, results []FailedRunItem, repo string, fr ghRunItem, idx int) {
	defer wg.Done()
	sem <- struct{}{}
	results[idx] = fetchAndBuildFailedRunItem(repo, fr)
	<-sem
}

func resolveFetchConcurrency(total int) int {
	if total <= 4 {
		return total
	}

	return 4
}

func initFailedRunTopLevel(p *PipelineErrorLogsPayload, fr ghRunItem) {
	p.Conclusion = fr.Conclusion
	if len(p.Conclusion) == 0 {
		p.Conclusion = "failure"
	}
	p.WorkflowName = fr.Name
	p.RunId = fr.DatabaseId
	p.Url = fr.Url
	p.Branch = fr.HeadBranch
	p.Sha = fr.HeadSha
	p.CreatedAt = fr.CreatedAt
	p.UpdatedAt = fr.UpdatedAt
	p.DurationSeconds = calculateRunDuration(fr.CreatedAt, fr.UpdatedAt)
	p.SavedLogFile = getCachedPipelineLogPath(fr.DatabaseId)
}

func fetchAndBuildFailedRunItem(repo string, fr ghRunItem) FailedRunItem {
	rawLogs := queryFailedRunLogs(repo, fr.DatabaseId)
	jobs := CorrelateRunFailedJobs(repo, fr.DatabaseId, rawLogs)
	item := buildBaseFailedRunItem(fr, rawLogs)
	item.FailedJobs = jobs
	item.StackTrace = extractStackTraceFromLog(rawLogs)

	return item
}

func buildBaseFailedRunItem(fr ghRunItem, rawLogs string) FailedRunItem {
	item := FailedRunItem{
		WorkflowName: fr.Name, RunId: fr.DatabaseId,
		Conclusion: fr.Conclusion, Status: fr.Status,
		Branch: fr.HeadBranch, Sha: fr.HeadSha,
		CreatedAt: fr.CreatedAt, UpdatedAt: fr.UpdatedAt,
		DurationSeconds: calculateRunDuration(fr.CreatedAt, fr.UpdatedAt),
		SavedLogFile:    getCachedPipelineLogPath(fr.DatabaseId),
		SavedMetaFile:   getCachedPipelineJSONPath(fr.DatabaseId),
		Url:             fr.Url, RawErrors: rawLogs,
	}

	return item
}

func collectFailedRuns(runs []ghRunItem) []ghRunItem {
	succeeded := make(map[string]bool)
	var activeFailed []ghRunItem

	for _, r := range runs {
		checkAndCollectRun(r, succeeded, &activeFailed)
	}

	return filterFailingRunsByTargetSha(activeFailed)
}

func normalizeWorkflowKey(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	lower = strings.TrimPrefix(lower, ".github/workflows/")
	lower = strings.TrimPrefix(lower, "workflows/")
	lower = strings.TrimSuffix(lower, ".yml")
	lower = strings.TrimSuffix(lower, ".yaml")

	return lower
}

func buildWorkflowScopeKey(r ghRunItem) string {
	nameKey := normalizeWorkflowKey(r.Name)
	branchKey := strings.ToLower(strings.TrimSpace(r.HeadBranch))
	shaKey := strings.ToLower(strings.TrimSpace(r.HeadSha))

	return branchKey + ":" + shaKey + ":" + nameKey
}

func checkAndCollectRun(r ghRunItem, succeeded map[string]bool, active *[]ghRunItem) {
	scopeKey := buildWorkflowScopeKey(r)
	if r.Conclusion == "success" {
		succeeded[scopeKey] = true

		return
	}
	if isFailingConclusion(r.Conclusion) && !succeeded[scopeKey] {
		succeeded[scopeKey] = true
		*active = append(*active, r)
	}
}

func filterFailingRunsByTargetSha(runs []ghRunItem) []ghRunItem {
	if len(runs) == 0 {
		return nil
	}

	return capFailedRuns(collectRunsMatchingSha(runs, runs[0].HeadSha), 5)
}

func collectRunsMatchingSha(runs []ghRunItem, targetSha string) []ghRunItem {
	var filtered []ghRunItem
	for _, r := range runs {
		if len(targetSha) == 0 || r.HeadSha == targetSha {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func capFailedRuns(runs []ghRunItem, limit int) []ghRunItem {
	if len(runs) > limit {
		return runs[:limit]
	}

	return runs
}

func buildLocalOrEmptyErrorPayload(payload PipelineErrorLogsPayload) PipelineErrorLogsPayload {
	localErr := readLocalLastErrorLog()

	if len(localErr) > 0 {
		payload.WorkflowName = "Local Command Execution"
		payload.Conclusion = "failure"
		payload.ErrorLogs = localErr

		return payload
	}

	payload.ErrorLogs = "No failed pipeline steps or local error logs found."

	return payload
}

func readLocalLastErrorLog() string {
	content, err := os.ReadFile(".gitmap/last_error.log")

	if err == nil && len(content) > 0 {
		return strings.TrimSpace(string(content))
	}

	return ""
}

func writeOrRenderErrorLogs(params ErrorLogOutputParams) error {
	contentToWrite, err := formatErrorLogContent(params)
	if err != nil {
		return err
	}
	_ = persistAutoErrorReport(params)

	if len(params.TempFile) > 0 || len(params.FilePath) > 0 {
		return writeErrorLogsToDisk(params, contentToWrite)
	}

	return dispatchErrorLogPresentation(params, contentToWrite)
}

func dispatchErrorLogPresentation(params ErrorLogOutputParams, content string) error {
	if params.IsJSON {
		return outputJSONErrorLogs(content)
	}
	if params.HasSuppressOutputLog {
		printSuppressedStagingNotice(params.Payload.SavedReportFile)

		return nil
	}
	renderErrorLogsTerminal(params.Payload)

	return nil
}

func outputJSONErrorLogs(content string) error {
	fmt.Println(content)
	_ = clipboard.WriteAll(content)

	return nil
}

func printSuppressedStagingNotice(reportFile string) {
	if len(reportFile) > 0 {
		fmt.Printf("  ✓ Error logs staged to %s (suppressed terminal output via --no-output-log)\n", reportFile)

		return
	}

	fmt.Println("  ✓ Error logs staged to filesystem (suppressed terminal output via --no-output-log)")
}

func writeErrorLogsToDisk(params ErrorLogOutputParams, content string) error {
	_ = clipboard.WriteAll(content)
	if len(params.TempFile) > 0 {
		targetPath := filepath.Join(resolveTempDir(), params.TempFile)

		return writeContentToFile(targetPath, content)
	}

	return writeContentToFile(params.FilePath, content)
}

func persistAutoErrorReport(params ErrorLogOutputParams) error {
	if !isFailingConclusion(params.Payload.Conclusion) && len(params.Payload.FailedRuns) == 0 {
		clearLocalErrorLogs()

		return nil
	}

	return saveActiveErrorReport(params.Payload)
}

func saveActiveErrorReport(p PipelineErrorLogsPayload) error {
	reportContent := p.ErrorLogs
	if len(reportContent) == 0 {
		reportContent = p.CombinedErrors
	}

	if len(reportContent) == 0 {
		return nil
	}

	_, err := writeCombinedErrorReport(reportContent)

	return err
}

func formatErrorLogContent(params ErrorLogOutputParams) (string, error) {
	if !params.IsJSON {
		return params.Payload.ErrorLogs, nil
	}

	b, err := json.MarshalIndent(params.Payload, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func renderErrorLogsTerminal(p PipelineErrorLogsPayload) {
	if p.IsRunning && len(p.FailedRuns) == 0 {
		renderActiveRunningBanner(p)
		copyReportToClipboard(buildClipboardCleanReport(p), false)

		return
	}

	if p.IsRunning {
		renderActiveRunningBanner(p)
	}

	renderAndCopyTerminalReport(p)
}

func renderAndCopyTerminalReport(p PipelineErrorLogsPayload) {
	hasFailure := isFailingConclusion(p.Conclusion) || len(p.FailedRuns) > 0
	if hasFailure {
		renderFailureTerminal(p)
		copyReportToClipboard(buildClipboardErrorReport(p), true)

		return
	}

	renderCleanSuccessTerminal(p)
	copyReportToClipboard(buildClipboardCleanReport(p), false)
}

func renderCleanSuccessTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Status: CLEAN (No errors found)%s\n",
		constants.ColorGreen, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Printf("\n  %s✓ All recent pipeline runs for %s on branch %s are PASSING (100%% green).%s\n\n",
		constants.ColorGreen, p.Repo, p.LatestBranch, constants.ColorReset)
	renderCleanSuccessDbAndHistory(p)
	printRerunETA(p.RerunEtaSeconds)
}

func renderPipelineMetaBlock(p PipelineErrorLogsPayload) {
	fmt.Printf("    %-18s %s\n", "Repo:", p.Repo)
	if len(p.RepoUrl) > 0 {
		fmt.Printf("    %-18s %s\n", "Repo URL:", p.RepoUrl)
	}
	if len(p.LatestBranch) > 0 {
		fmt.Printf("    %-18s %s\n", "Branch:", p.LatestBranch)
	}
	if len(p.LastHash) > 0 {
		fmt.Printf("    %-18s %s\n", "Last Commit:", p.LastHash)
	}
	renderPipelineMetaVersionAndPR(p)
}

func renderPipelineMetaVersionAndPR(p PipelineErrorLogsPayload) {
	if len(p.LastReleaseVersion) > 0 {
		fmt.Printf("    %-18s %s\n", "Last Release:", p.LastReleaseVersion)
	}
	fmt.Printf("    %-18s %d\n", "Open PRs:", p.OpenPRsCount)
	if len(p.Status) > 0 && len(p.Conclusion) > 0 {
		fmt.Printf("    %-18s %s (conclusion: %s)\n", "Run Status:", p.Status, p.Conclusion)
	}
}

func renderCleanSuccessDbAndHistory(p PipelineErrorLogsPayload) {
	if len(p.DbPath) > 0 {
		fmt.Printf("  • Pipeline DB:     %s\n", FormatRelativeDbPath(p.DbPath))
		fmt.Printf("  • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
	}

	runs := queryWorkflowRuns(p.Repo)
	RenderHistorySummaryTable(runs)
}

func renderActiveRunningBanner(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Active Pipeline is RUNNING%s: [%s #%d] (ETA: %ds)\n",
		constants.ColorYellow, constants.ColorReset, p.ActiveRunName, p.ActiveRunId, p.EtaSeconds)
	if len(p.ActiveRunUrl) > 0 {
		fmt.Printf("    URL: %s\n\n", p.ActiveRunUrl)
	}
}

func renderFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Failure Detected%s\n",
		constants.ColorRed, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Println()
	if len(p.FailedRuns) == 0 {
		renderSingleFailureTerminal(p)

		return
	}

	renderFailureSectionsAndETA(p)
}

func renderFailureSectionsAndETA(p PipelineErrorLogsPayload) {
	renderCombinedSectionsTerminal(p.SectionFailures)
	renderFailedRunsBreakdown(p.FailedRuns)
	renderSavedLocationsTerminal(p)
	runs := queryWorkflowRuns(p.Repo)
	RenderHistorySummaryTable(runs)
	printRerunETA(p.RerunEtaSeconds)
}

func buildClipboardErrorReport(p PipelineErrorLogsPayload) string {
	var sb strings.Builder
	appendClipboardMetaHeader(&sb, "GITMAP PIPELINE ERROR REPORT", p)
	appendClipboardBodyContent(&sb, p)

	return strings.TrimSpace(sb.String())
}

func appendClipboardBodyContent(sb *strings.Builder, p PipelineErrorLogsPayload) {
	if len(p.ErrorLogs) > 0 {
		sb.WriteString(p.ErrorLogs)
		sb.WriteString("\n")

		return
	}
	if len(p.CombinedErrors) > 0 {
		sb.WriteString(p.CombinedErrors)
		sb.WriteString("\n")
	}
}

func buildClipboardCleanReport(p PipelineErrorLogsPayload) string {
	var sb strings.Builder
	appendClipboardMetaHeader(&sb, "GITMAP PIPELINE STATUS: CLEAN (No errors found)", p)
	sb.WriteString("All recent pipeline workflow runs are PASSING (100% green).\n")

	return strings.TrimSpace(sb.String())
}

func appendClipboardMetaHeader(sb *strings.Builder, title string, p PipelineErrorLogsPayload) {
	sb.WriteString("================================================================================\n")
	sb.WriteString(title + "\n")
	sb.WriteString("================================================================================\n")
	sb.WriteString(fmt.Sprintf("Repo:                 %s\n", p.Repo))
	if len(p.RepoUrl) > 0 {
		sb.WriteString(fmt.Sprintf("Repo URL:             %s\n", p.RepoUrl))
	}
	appendClipboardMetaDetails(sb, p)
}

func appendClipboardMetaDetails(sb *strings.Builder, p PipelineErrorLogsPayload) {
	sb.WriteString(fmt.Sprintf("Branch:               %s\n", p.LatestBranch))
	sb.WriteString(fmt.Sprintf("Last Commit:          %s\n", p.LastHash))
	sb.WriteString(fmt.Sprintf("Last Release:         %s\n", p.LastReleaseVersion))
	sb.WriteString(fmt.Sprintf("Open PRs:             %d\n", p.OpenPRsCount))
	if len(p.Status) > 0 && len(p.Conclusion) > 0 {
		sb.WriteString(fmt.Sprintf("Status:               %s (conclusion: %s)\n", p.Status, p.Conclusion))
	}
	if len(p.Url) > 0 {
		sb.WriteString(fmt.Sprintf("Pipeline Run URL:     %s\n", p.Url))
	}
	sb.WriteString("================================================================================\n\n")
}

func copyReportToClipboard(content string, isFailure bool) {
	if len(content) == 0 {
		return
	}

	err := clipboard.WriteAll(content)
	if err != nil {
		return
	}

	printClipboardNotice(isFailure)
}

func printClipboardNotice(isFailure bool) {
	if isFailure {
		fmt.Printf("\n  📋 Copied pipeline error logs to clipboard\n")

		return
	}

	fmt.Printf("\n  📋 Copied pipeline status to clipboard\n")
}

func renderCombinedSectionsTerminal(sections []SectionFailure) {
	if len(sections) == 0 {
		return
	}

	fmt.Printf("  %s● Combined Pipeline Section Failures [%d failed section(s)]:%s\n",
		constants.ColorRed, len(sections), constants.ColorReset)
	for i, sec := range sections {
		renderSingleSectionFailureRow(sec, i+1, len(sections))
	}

	fmt.Println()
}

func renderSingleSectionFailureRow(sec SectionFailure, idx, total int) {
	fmt.Printf("    %s[%d/%d] %s #%d ➔ Job: %s | Step: %s%s\n",
		constants.ColorCyan, idx, total, sec.WorkflowName, sec.RunId, sec.JobName, sec.StepName, constants.ColorReset)
	if len(sec.FailureSummary) > 0 {
		fmt.Printf("      Error:   %s%s%s\n", constants.ColorRed, sec.FailureSummary, constants.ColorReset)
	}

	renderSectionErrorLines(sec.ErrorLines)
	renderSectionStackTrace(sec.StackTrace)

	if len(sec.SavedLogFile) > 0 {
		fmt.Printf("      Log:     %s\n", sec.SavedLogFile)
	}
}

func renderSectionStackTrace(stack string) {
	if len(stack) == 0 {
		return
	}

	fmt.Printf("      Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		fmt.Printf("        %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
	}
}

func renderSectionErrorLines(lines []string) {
	if len(lines) == 0 {
		return
	}

	fmt.Printf("      Details:\n")
	capped := capErrorLines(lines, 12)
	for _, l := range capped {
		fmt.Printf("        %s%s%s\n", constants.ColorYellow, l, constants.ColorReset)
	}

	printRemainingLineCount(len(lines), len(capped))
}

func printRemainingLineCount(total, capped int) {
	if total > capped {
		fmt.Printf("        %s... (%d more lines in log)%s\n", constants.ColorDim, total-capped, constants.ColorReset)
	}
}

func renderFailedRunsBreakdown(failedRuns []FailedRunItem) {
	total := len(failedRuns)
	fmt.Printf("  %s● Detailed Failure Logs Across Workflow Runs [%d run(s)]:%s\n\n",
		constants.ColorRed, total, constants.ColorReset)
	for i, fr := range failedRuns {
		renderFailedRunCard(fr, i+1, total)
	}
}

func renderSingleFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Latest Pipeline Failure [%s #%d]:%s\n\n",
		constants.ColorRed, p.WorkflowName, p.RunId, constants.ColorReset)
	clean := extractCleanErrorLines(p.ErrorLogs)
	printLogsContent(clean, p.ErrorLogs)
	renderSavedLocationsTerminal(p)
	printRerunETA(p.RerunEtaSeconds)
}

func renderFailedRunCard(fr FailedRunItem, idx, total int) {
	fmt.Printf("  %s┌─ [%d/%d] %s #%d ──────────────────────────%s\n",
		constants.ColorRed, idx, total, fr.WorkflowName, fr.RunId, constants.ColorReset)
	renderRunCardMeta(fr)
	for _, job := range fr.FailedJobs {
		renderFailedJobSection(job)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────%s\n\n",
		constants.ColorRed, constants.ColorReset)
}

func renderRunCardMeta(fr FailedRunItem) {
	if len(fr.CreatedAt) > 0 {
		fmt.Printf("  │ When Run:  %s\n", formatRunTimestamp(fr.CreatedAt))
	}
	if fr.DurationSeconds > 0 {
		fmt.Printf("  │ Duration:  %s\n", formatDurationSeconds(fr.DurationSeconds))
	}
	renderRunCardBranchAndLog(fr)
}

func renderRunCardBranchAndLog(fr FailedRunItem) {
	if len(fr.Branch) > 0 {
		fmt.Printf("  │ Branch:    %s | Commit: %s\n", fr.Branch, fr.Sha)
	}
	if len(fr.SavedLogFile) > 0 {
		fmt.Printf("  │ Saved Log: %s\n", fr.SavedLogFile)
	}
	if len(fr.Url) > 0 {
		fmt.Printf("  │ URL:       %s\n", fr.Url)
	}
}

func renderSavedLocationsTerminal(p PipelineErrorLogsPayload) {
	fmt.Println("  💾 Pipeline Error Logs & Artifacts:")
	if len(p.SavedReportFile) > 0 {
		fmt.Printf("    • Combined Report: %s\n", p.SavedReportFile)
	}
	if len(p.SavedLogFile) > 0 {
		fmt.Printf("    • Latest Run Log:  %s\n", p.SavedLogFile)
	}
	renderSavedDbAndUrl(p)
}

func renderSavedDbAndUrl(p PipelineErrorLogsPayload) {
	if len(p.DbPath) > 0 {
		fmt.Printf("    • Pipeline DB:     %s\n", FormatRelativeDbPath(p.DbPath))
		fmt.Printf("    • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
	}
	if len(p.Url) > 0 {
		fmt.Printf("    • Web Run URL:     %s\n\n", p.Url)
	}
}

func renderFailedJobSection(job FailedJobItem) {
	fmt.Printf("  │ Job:   %s\n", job.JobName)
	fmt.Printf("  │ Step:  %s\n", job.StepName)
	if len(job.FailureSummary) > 0 {
		fmt.Printf("  │ Error: %s%s%s\n", constants.ColorRed, job.FailureSummary, constants.ColorReset)
	}

	for _, line := range job.ErrorLines {
		fmt.Printf("  │   %s\n", line)
	}
	renderJobStackTrace(job.StackTrace)
}

func renderJobStackTrace(stack string) {
	if len(stack) == 0 {
		return
	}

	fmt.Printf("  │ Stack Trace:\n")
	for _, line := range strings.Split(strings.TrimSpace(stack), "\n") {
		fmt.Printf("  │   %s%s%s\n", constants.ColorDim, line, constants.ColorReset)
	}
}

func printLogsContent(clean, raw string) {
	if len(clean) > 0 {
		fmt.Println(clean)

		return
	}

	fmt.Println(raw)
}

func printRerunETA(eta int) {
	if eta <= 0 {
		return
	}

	fmt.Printf("\n  %s● Estimated pipeline rerun duration (ETA): ~%ds%s\n",
		constants.ColorYellow, eta, constants.ColorReset)
	fmt.Println("    (Based on historical successful pipeline runs baseline)")
}

func printPipelineErrorLogsHelp() {
	fmt.Println("Usage: gitmap pipeline error-logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --timeline          Watch pipeline dynamic timeline until completion")
	fmt.Println("  -f, --fix               Execute internal CI/CD diagnostic & auto-repair suite")
	fmt.Println("  -c, --check             Run internal CI/CD checks without modifying files")
	fmt.Println("  -v, --detailed, --verbose  Show full raw error logs including passing ok lines")
	fmt.Println("  -y, --yes               Auto-confirm prompts non-interactively")
	fmt.Println("  --json                  Output data in structured JSON format")
	fmt.Println("  --file <path>           Write error logs to specified file path")
	fmt.Println("  --tempfile <filename>   Write error logs to .ai-memory/temp/<filename>")
	fmt.Println("  -n, --no-output-log     Stage error logs to disk without displaying in terminal")
}

func printPipelineLogsHelp() {
	fmt.Println("Usage: gitmap pipeline logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --json                  Output workflow status and URL in JSON format")
}
