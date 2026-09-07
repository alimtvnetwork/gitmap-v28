package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	payload := PipelineErrorLogsPayload{
		Repo: repo,
	}

	if len(runs) == 0 {
		return buildLocalOrEmptyErrorPayload(payload)
	}

	latest := runs[0]
	payload.WorkflowName = latest.Name
	payload.RunId = latest.DatabaseId
	payload.Status = latest.Status
	payload.Conclusion = latest.Conclusion
	payload.Url = latest.Url

	isRunning := latest.Status == "in_progress" || latest.Status == "queued"
	if isRunning {
		payload.IsRunning = true
		payload.EtaSeconds = calculateETA(runs)
		payload.ErrorLogs = fmt.Sprintf("Pipeline is currently running. Estimated completion in %d seconds.", payload.EtaSeconds)
	}

	failedRuns := collectFailedRuns(runs)
	if len(failedRuns) > 0 {
		populateFailedRunsPayload(repo, failedRuns, &payload)

		return payload
	}

	if isRunning {
		return payload
	}

	return buildLocalOrEmptyErrorPayload(payload)
}

func populateFailedRunsPayload(repo string, failedRuns []ghRunItem, p *PipelineErrorLogsPayload) {
	p.Conclusion = "failure"
	p.WorkflowName = failedRuns[0].Name
	p.RunId = failedRuns[0].DatabaseId
	p.Url = failedRuns[0].Url

	for _, fr := range failedRuns {
		p.FailedRuns = append(p.FailedRuns, fetchAndBuildFailedRunItem(repo, fr))
	}

	p.ErrorLogs = formatAggregatedErrorLogs(p.FailedRuns)
}

func fetchAndBuildFailedRunItem(repo string, fr ghRunItem) FailedRunItem {
	rawLogs := queryFailedRunLogs(repo, fr.DatabaseId)
	jobs := ParseFailedLogLines(rawLogs)

	return FailedRunItem{
		WorkflowName: fr.Name,
		RunId:        fr.DatabaseId,
		Conclusion:   fr.Conclusion,
		Url:          fr.Url,
		FailedJobs:   jobs,
		RawErrors:    rawLogs,
	}
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

	if len(params.TempFile) > 0 {
		tempDir := resolveTempDir()
		targetPath := filepath.Join(tempDir, params.TempFile)

		return writeContentToFile(targetPath, contentToWrite)
	}

	if len(params.FilePath) > 0 {
		return writeContentToFile(params.FilePath, contentToWrite)
	}

	if params.IsJSON {
		fmt.Println(contentToWrite)

		return nil
	}

	renderErrorLogsTerminal(params.Payload)

	return nil
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
	if p.IsRunning {
		fmt.Printf("  %s● Pipeline is RUNNING%s (ETA: %ds)\n", constants.ColorYellow, constants.ColorReset, p.EtaSeconds)

		return
	}

	if p.Conclusion == "failure" {
		renderFailureTerminal(p)

		return
	}

	fmt.Printf("  %s● No error logs found.%s Status: %s (conclusion: %s)\n",
		constants.ColorGreen, constants.ColorReset, p.Status, p.Conclusion)
	printRerunETA(p.RerunEtaSeconds)
}

func renderFailureTerminal(p PipelineErrorLogsPayload) {
	if len(p.FailedRuns) == 0 {
		renderSingleFailureTerminal(p)

		return
	}

	total := len(p.FailedRuns)
	fmt.Printf("  %s● Pipeline Failures [%d failed workflow run(s)]:%s\n\n",
		constants.ColorRed, total, constants.ColorReset)

	for i, fr := range p.FailedRuns {
		renderFailedRunCard(fr, i+1, total)
	}

	printRerunETA(p.RerunEtaSeconds)
}

func renderSingleFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Latest Pipeline Failure [%s #%d]:%s\n\n",
		constants.ColorRed, p.WorkflowName, p.RunId, constants.ColorReset)
	clean := extractCleanErrorLines(p.ErrorLogs)
	printLogsContent(clean, p.ErrorLogs)
	printRerunETA(p.RerunEtaSeconds)
}

func renderFailedRunCard(fr FailedRunItem, idx, total int) {
	fmt.Printf("  %s┌─ [%d/%d] %s #%d ──────────────────────────%s\n",
		constants.ColorRed, idx, total, fr.WorkflowName, fr.RunId, constants.ColorReset)
	for _, job := range fr.FailedJobs {
		renderFailedJobSection(job)
	}
	if len(fr.Url) > 0 {
		fmt.Printf("  │ URL:   %s\n", fr.Url)
	}
	fmt.Printf("  %s└──────────────────────────────────────────────────────────%s\n\n",
		constants.ColorRed, constants.ColorReset)
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
