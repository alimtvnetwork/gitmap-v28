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
	isJSON := hasArgFlag(args, "--json")
	filePath := extractFlagVal(args, "--file")
	tempFileName := extractFlagVal(args, "--tempfile")

	repo := resolveCurrentRepoSlug()
	runs := queryWorkflowRuns(repo)
	payload := buildErrorLogsPayload(repo, runs)

	return writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:  payload,
		IsJSON:   isJSON,
		FilePath: filePath,
		TempFile: tempFileName,
	})
}

func handlePipelineLogs(args []string) error {
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

	rawLogs := queryAllRunLogs(repo, latestRun.DatabaseID)

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
	payload.RunID = latest.DatabaseID
	payload.Status = latest.Status
	payload.Conclusion = latest.Conclusion
	payload.URL = latest.URL

	isRunning := latest.Status == "in_progress" || latest.Status == "queued"

	if isRunning {
		payload.IsRunning = true
		payload.EtaSeconds = calculateETA(runs)
		payload.ErrorLogs = fmt.Sprintf("Pipeline is currently running. Estimated completion in %d seconds.", payload.EtaSeconds)

		return payload
	}

	targetRunID := findFailedRunID(runs, &payload)

	if targetRunID > 0 {
		payload.ErrorLogs = queryFailedRunLogs(repo, targetRunID)

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

func findFailedRunID(runs []ghRunItem, p *PipelineErrorLogsPayload) int64 {
	for _, r := range runs {
		if r.Conclusion == "failure" {
			p.WorkflowName = r.Name
			p.RunID = r.DatabaseID
			p.Conclusion = r.Conclusion
			p.URL = r.URL

			return r.DatabaseID
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
			constants.ColorRed, p.WorkflowName, p.RunID, constants.ColorReset)
		fmt.Println(p.ErrorLogs)

		return
	}

	fmt.Printf("  %s● No error logs found.%s Status: %s (conclusion: %s)\n",
		constants.ColorGreen, constants.ColorReset, p.Status, p.Conclusion)
}
