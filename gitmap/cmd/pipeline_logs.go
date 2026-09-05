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

	maybeOfferAutoFix(payload, isJSON, wantFix, wantCheck, filePath, tempFileName, args)

	return nil
}

func maybeOfferAutoFix(p PipelineErrorLogsPayload, isJSON, wantFix, wantCheck bool, file, temp string, args []string) {
	isTargetingFile := len(file) > 0 || len(temp) > 0
	if isJSON || wantFix || wantCheck || isTargetingFile {
		return
	}
	if p.Conclusion != "failure" {
		return
	}
	if confirmOrSkip("Would you like to run internal CI/CD diagnostic & auto-repair scripts?", args) {
		runInternalCICDChecks(true)
	}
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

		return payload
	}

	targetRunId := findFailedRunId(runs, &payload)

	if targetRunId > 0 {
		payload.ErrorLogs = queryFailedRunLogs(repo, targetRunId)

		return payload
	}

	return buildLocalOrEmptyErrorPayload(payload)
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

func findFailedRunId(runs []ghRunItem, p *PipelineErrorLogsPayload) uint64 {
	for _, r := range runs {
		if r.Conclusion == "failure" {
			p.WorkflowName = r.Name
			p.RunId = r.DatabaseId
			p.Conclusion = r.Conclusion
			p.Url = r.Url

			return r.DatabaseId
		}
	}

	return 0
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
		fmt.Printf("  %s● Latest Pipeline Failure [%s #%d]:%s\n\n",
			constants.ColorRed, p.WorkflowName, p.RunId, constants.ColorReset)
		fmt.Println(p.ErrorLogs)
		printRerunETA(p.RerunEtaSeconds)

		return
	}

	fmt.Printf("  %s● No error logs found.%s Status: %s (conclusion: %s)\n",
		constants.ColorGreen, constants.ColorReset, p.Status, p.Conclusion)
	printRerunETA(p.RerunEtaSeconds)
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
