package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func resolvePipelineDir() string {
	if configuredDir := getSettingPipelineDir(); configuredDir != "" {
		return configuredDir
	}

	return filepath.Join(resolveRepoRootDir(), ".gitmap", "pipeline")
}

func getSettingPipelineDir() string {
	db, err := openDB()
	if err != nil {
		return ""
	}
	defer db.Close()

	return db.GetSetting("pipeline.dir")
}

func resolveRepoRootDir() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	return "."
}

func readCachedPipelineLog(runId uint64) (string, bool) {
	dir := resolvePipelineDir()
	filePath := filepath.Join(dir, fmt.Sprintf("%d.log", runId))
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return "", false
	}

	return string(data), true
}

func writeCachedPipelineLog(runId uint64, logContent, repo string) error {
	dir := resolvePipelineDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(dir, fmt.Sprintf("%d.log", runId))
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		return err
	}

	return writeCachedPipelineJSON(dir, runId, logFile, repo, len(logContent))
}

func writeCachedPipelineJSON(dir string, runId uint64, logFile, repo string, byteCount int) error {
	meta := map[string]any{
		"runId":     runId,
		"repo":      repo,
		"cachedAt":  time.Now().UTC().Format(time.RFC3339),
		"logFile":   logFile,
		"byteCount": byteCount,
	}

	jsonBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	jsonFile := filepath.Join(dir, fmt.Sprintf("%d.json", runId))

	return os.WriteFile(jsonFile, jsonBytes, 0644)
}

func getCachedPipelineLogPath(runId uint64) string {
	return filepath.Join(resolvePipelineDir(), fmt.Sprintf("%d.log", runId))
}

func getCachedPipelineJSONPath(runId uint64) string {
	return filepath.Join(resolvePipelineDir(), fmt.Sprintf("%d.json", runId))
}

func resolvePipelineErrorReportPath() string {
	return filepath.Join(resolvePipelineDir(), "pipeline_errors.log")
}

func writeCombinedErrorReport(content string) (string, error) {
	dir := resolvePipelineDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	reportPath := resolvePipelineErrorReportPath()
	if err := os.WriteFile(reportPath, []byte(content), 0644); err != nil {
		return "", err
	}

	_ = writeLastErrorLog(content)

	return reportPath, nil
}

func writeLastErrorLog(content string) error {
	gitmapDir := filepath.Join(resolveRepoRootDir(), ".gitmap")
	_ = os.MkdirAll(gitmapDir, 0755)
	lastErrFile := filepath.Join(gitmapDir, "last_error.log")

	return os.WriteFile(lastErrFile, []byte(content), 0644)
}
