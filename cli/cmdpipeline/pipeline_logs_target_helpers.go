package cmdpipeline

import (
	"encoding/json"
	"os"
	"strings"
)

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
	effectiveLimit := resolveFailedRunsLimit(limit)
	if len(runs) > effectiveLimit {
		return runs[:effectiveLimit]
	}

	return runs
}

func resolveFailedRunsLimit(limit int) int {
	if activeLogLineLimit > 0 {
		return activeLogLineLimit
	}

	return limit
}

func buildLocalOrEmptyErrorPayload(payload PipelineErrorLogsPayload) PipelineErrorLogsPayload {
	localErr := readLocalLastErrorLog()

	if len(localErr) > 0 {
		payload.WorkflowName = "Local Command Execution"
		payload.Conclusion = "failure"
		payload.ErrorLogs = localErr

		return payload
	}

	payload.WorkflowName = "CI/CD Pipeline"
	payload.Conclusion = "success"
	payload.ErrorLogs = "No failed pipeline steps or local error logs found."

	return payload
}

func readLocalLastErrorLog() string {
	content, err := os.ReadFile(".gitmap/last_error.log")
	if err != nil || len(content) == 0 {
		return ""
	}

	trimmed := strings.TrimSpace(string(content))
	var report struct {
		Command string `json:"command"`
		Error   string `json:"error"`
	}
	if jsonErr := json.Unmarshal([]byte(trimmed), &report); jsonErr != nil {
		return trimmed
	}
	if report.Command != "" && !isPipelineCommandName(report.Command) {
		return ""
	}
	if len(report.Error) > 0 {
		return report.Error
	}

	return trimmed
}

func isPipelineCommandName(cmd string) bool {
	c := strings.ToLower(strings.TrimSpace(cmd))
	return c == "pipeline" || c == "pipeline-ai" || c == "pl" || c == "ci" || c == "pe" || c == "te" ||
		strings.HasPrefix(c, "pipeline") || strings.HasPrefix(c, "ci")
}
