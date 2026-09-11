package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
)

func handlePipelineErrorLogs(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineErrorLogsHelp()
		return nil
	}

	isJSON := hasArgFlag(args, "--json")
	hasTimeout := hasArgFlag(args, "-t") || hasArgFlag(args, "--timeout") || hasArgFlag(args, "--timeline") || hasArgFlag(args, "-w") || hasArgFlag(args, "--watch")
	wantFix := hasArgFlag(args, "--fix") || hasArgFlag(args, "-f")
	wantCheck := hasArgFlag(args, "--check") || hasArgFlag(args, "-c")
	filePath := extractFlagVal(args, "--file")
	tempFileName := extractFlagVal(args, "--tempfile")

	repo := resolveCurrentRepoSlug()

	if hasTimeout {
		return runPipelineErrorLogsDynamicTimeline(ErrorLogsTimelineParams{
			Repo:         repo,
			IsJSON:       isJSON,
			WantFix:      wantFix,
			WantCheck:    wantCheck,
			FilePath:     filePath,
			TempFileName: tempFileName,
			Args:         args,
		})
	}

	runs := queryWorkflowRuns(repo)
	payload := buildErrorLogsPayload(repo, runs)
	if len(runs) > 0 {
		payload.RerunEtaSeconds = calculateAverageDuration(runs, payload.WorkflowName)
	}

	if wantFix || wantCheck {
		payload.CICDChecks = runInternalCICDChecks(wantFix)
	}

	err := writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:  payload,
		IsJSON:   isJSON,
		FilePath: filePath,
		TempFile: tempFileName,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlePipelineLogs(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") {
		printPipelineLogsHelp()
		return nil
	}

	repo := resolveCurrentRepoSlug()
	runs := queryWorkflowRuns(repo)

	if len(runs) == 0 {
		fmt.Printf("No recent pipeline runs found for %s.\n", repo)

		return nil
	}

	latestRun := runs[0]
	isRunning := latestRun.Status == "in_progress" || latestRun.Status == "queued"

	if isRunning {
		eta := calculateETA(runs)
		fmt.Printf("● Pipeline [%s] is currently RUNNING (ETA: %ds)\n", latestRun.Name, eta)
	}

	rawLogs := queryAllRunLogs(repo, latestRun.DatabaseId)

	if len(rawLogs) > 0 {
		fmt.Println(rawLogs)
	} else {
		fmt.Printf("● Pipeline [%s] status: %s (conclusion: %s)\n", latestRun.Name, latestRun.Status, latestRun.Conclusion)
	}

	return nil
}

func buildErrorLogsPayload(repo string, runs []ghRunItem) PipelineErrorLogsPayload {
	payload := initBaseErrorLogsPayload(repo)
	if len(runs) == 0 {
		return buildLocalOrEmptyErrorPayload(payload)
	}
	initLatestRunMeta(&payload, runs[0])
	checkAndApplyRunningState(&payload, runs)
	failedRuns := resolveFailedRunsForPayload(repo, runs)
	if len(failedRuns) > 0 {
		populateFailedRunsPayload(repo, failedRuns, &payload)

		return payload
	}
	if payload.IsRunning {
		return payload
	}

	return buildLocalOrEmptyErrorPayload(payload)
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
	failed := collectFailedRuns(runs)
	if len(failed) == 0 {
		return queryRecentFailedRuns(repo, 5)
	}

	return failed
}

func populateFailedRunsPayload(repo string, failedRuns []ghRunItem, p *PipelineErrorLogsPayload) {
	initFailedRunTopLevel(p, failedRuns[0])
	for _, fr := range failedRuns {
		p.FailedRuns = append(p.FailedRuns, fetchAndBuildFailedRunItem(repo, fr))
	}
	p.SectionFailures = extractAllSectionFailures(p.FailedRuns)
	p.CombinedErrors = formatCombinedSectionFailures(p.SectionFailures)
	p.ErrorLogs = formatAggregatedErrorLogs(p.FailedRuns)
}

func initFailedRunTopLevel(p *PipelineErrorLogsPayload, fr ghRunItem) {
	p.Conclusion = "failure"
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
	jobs := ParseFailedLogLines(rawLogs)
	item := buildBaseFailedRunItem(fr, rawLogs)
	item.FailedJobs = jobs

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
	var targetSha string
	var failed []ghRunItem

	for _, r := range runs {
		if r.Conclusion != "failure" {
			continue
		}
		if len(targetSha) == 0 && len(r.HeadSha) > 0 {
			targetSha = r.HeadSha
		}
		if len(targetSha) > 0 && r.HeadSha == targetSha {
			failed = append(failed, r)
			continue
		}
		if len(targetSha) == 0 {
			failed = append(failed, r)
		}
	}

	if len(failed) > 5 {
		return failed[:5]
	}

	return failed
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
	if params.IsJSON {
		fmt.Println(contentToWrite)

		return nil
	}
	renderErrorLogsTerminal(params.Payload)

	return nil
}

func writeErrorLogsToDisk(params ErrorLogOutputParams, content string) error {
	if len(params.TempFile) > 0 {
		targetPath := filepath.Join(resolveTempDir(), params.TempFile)

		return writeContentToFile(targetPath, content)
	}

	return writeContentToFile(params.FilePath, content)
}

func persistAutoErrorReport(params ErrorLogOutputParams) error {
	if params.Payload.Conclusion != "failure" && len(params.Payload.FailedRuns) == 0 {
		return nil
	}

	reportContent := params.Payload.ErrorLogs
	if len(reportContent) == 0 {
		reportContent = params.Payload.CombinedErrors
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

		return
	}
	if p.IsRunning {
		renderActiveRunningBanner(p)
	}

	if p.Conclusion == "failure" || len(p.FailedRuns) > 0 {
		renderFailureTerminal(p)

		return
	}

	fmt.Printf("  %s● No error logs found.%s Status: %s (conclusion: %s)\n",
		constants.ColorGreen, constants.ColorReset, p.Status, p.Conclusion)
	printRerunETA(p.RerunEtaSeconds)
}

func renderActiveRunningBanner(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Active Pipeline is RUNNING%s: [%s #%d] (ETA: %ds)\n",
		constants.ColorYellow, constants.ColorReset, p.ActiveRunName, p.ActiveRunId, p.EtaSeconds)
	if len(p.ActiveRunUrl) > 0 {
		fmt.Printf("    URL: %s\n\n", p.ActiveRunUrl)
	}
}

func renderFailureTerminal(p PipelineErrorLogsPayload) {
	if len(p.FailedRuns) == 0 {
		renderSingleFailureTerminal(p)

		return
	}

	renderCombinedSectionsTerminal(p.SectionFailures)
	renderFailedRunsBreakdown(p.FailedRuns)
	renderSavedLocationsTerminal(p)
	printRerunETA(p.RerunEtaSeconds)
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
		fmt.Printf("      Error: %s%s%s\n", constants.ColorRed, sec.FailureSummary, constants.ColorReset)
	}
	if len(sec.SavedLogFile) > 0 {
		fmt.Printf("      Log:   %s\n", sec.SavedLogFile)
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
	if len(p.DbPath) > 0 {
		fmt.Printf("    • Pipeline DB:     %s\n", p.DbPath)
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
	fmt.Println("  -y, --yes               Auto-confirm prompts non-interactively")
	fmt.Println("  --json                  Output data in structured JSON format")
	fmt.Println("  --file <path>           Write error logs to specified file path")
	fmt.Println("  --tempfile <filename>   Write error logs to .lovable/temp/<filename>")
}

func printPipelineLogsHelp() {
	fmt.Println("Usage: gitmap pipeline logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --json                  Output workflow status and URL in JSON format")
}
