package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func resolvePipelineDir() string {
	if configuredDir := getSettingPipelineDir(); configuredDir != "" {
		return configuredDir
	}

	return pipelinedb.PipelineDbDir()
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

func resolvePipelineDirForRepo(repo string) string {
	targetRepo := repo
	if len(targetRepo) == 0 {
		targetRepo = resolveCurrentRepoSlug()
	}

	return pipelinedb.RepoPipelineDir(targetRepo)
}

func isCorruptOrFallbackErrorLog(content string) bool {
	if strings.Contains(content, "gh command failed") {
		return true
	}
	if strings.Contains(content, "at github.com/alimtvnetwork/") {
		return true
	}
	if strings.Contains(content, "Unable to fetch failed logs via gh CLI") {
		return true
	}

	return false
}

func readCachedPipelineLog(runId uint64) (string, bool) {
	return readCachedPipelineLogForRepo("", runId)
}

func readCachedPipelineLogForRepo(repo string, runId uint64) (string, bool) {
	candidates := []string{
		filepath.Join(resolvePipelineDirForRepo(repo), fmt.Sprintf("%d.log", runId)),
		filepath.Join(resolvePipelineDir(), fmt.Sprintf("%d.log", runId)),
	}
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		content := string(data)
		hasValidData := err == nil && len(data) > 0 && !isCorruptOrFallbackErrorLog(content)
		if hasValidData {
			return content, true
		}
	}

	return readCachedLogFromDb(repo, runId)
}

func readCachedLogFromDb(repo string, runId uint64) (string, bool) {
	pipeDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return "", false
	}
	defer pipeDb.Close()

	if content, ok := queryDetailLogFromDb(pipeDb, runId); ok && !isCorruptOrFallbackErrorLog(content) {
		return content, true
	}
	if content, ok := queryCompactLogFromDb(pipeDb, runId); ok && !isCorruptOrFallbackErrorLog(content) {
		return content, true
	}
	if content, ok := queryLegacyErrorLogFromDb(pipeDb, runId); ok && !isCorruptOrFallbackErrorLog(content) {
		return content, true
	}

	return "", false
}

func queryDetailLogFromDb(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) (string, bool) {
	detailRes := pipeDb.QueryDetailedErrorLogsByRunId(runId)
	if !detailRes.HasRecord() {
		return "", false
	}
	for _, rec := range detailRes.Data {
		if len(rec.RawLogs) > 0 {
			return rec.RawLogs, true
		}
		if len(rec.ErrorText) > 0 {
			return rec.ErrorText, true
		}
	}

	return "", false
}

func queryCompactLogFromDb(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) (string, bool) {
	compactRes := pipeDb.QueryCompactErrorLogsByRunId(runId)
	if !compactRes.HasRecord() {
		return "", false
	}
	for _, rec := range compactRes.Data {
		if len(rec.CompactLogs) > 0 {
			return rec.CompactLogs, true
		}
		if len(rec.ErrorText) > 0 {
			return rec.ErrorText, true
		}
	}

	return "", false
}

func queryLegacyErrorLogFromDb(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) (string, bool) {
	legacyRes := pipeDb.QueryErrorLogsByRunId(runId)
	if legacyRes.HasRecord() && len(legacyRes.Data[0].RawLogs) > 0 {
		return legacyRes.Data[0].RawLogs, true
	}

	return "", false
}

func writeCachedPipelineLog(runId uint64, logContent, repo string) error {
	dir := resolvePipelineDirForRepo(repo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	logFile := filepath.ToSlash(filepath.Join(dir, fmt.Sprintf("%d.log", runId)))
	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		return err
	}

	persistLogToRepoSplitDb(repo, runId, logContent)

	return writeCachedPipelineJSON(dir, runId, logFile, repo, len(logContent))
}

func readCachedPipelineJobs(runId uint64, repo string) ([]ghJobItem, bool) {
	candidates := []string{
		filepath.Join(resolvePipelineDirForRepo(repo), fmt.Sprintf("%d.jobs.json", runId)),
		filepath.Join(resolvePipelineDir(), fmt.Sprintf("%d.jobs.json", runId)),
	}
	for _, p := range candidates {
		if jobs, isLoaded := loadJobsFromFile(p); isLoaded {
			return jobs, true
		}
	}

	return nil, false
}

func loadJobsFromFile(path string) ([]ghJobItem, bool) {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return nil, false
	}

	var jobs []ghJobItem
	if err := json.Unmarshal(data, &jobs); err != nil || len(jobs) == 0 {
		return nil, false
	}

	return jobs, true
}

func writeCachedPipelineJobs(runId uint64, repo string, jobs []ghJobItem) error {
	if len(jobs) == 0 {
		return nil
	}

	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}

	repoDir := resolvePipelineDirForRepo(repo)
	_ = os.MkdirAll(repoDir, 0755)

	return os.WriteFile(filepath.Join(repoDir, fmt.Sprintf("%d.jobs.json", runId)), data, 0644)
}

func buildPersistRecords(repo string, runId uint64, workflow, raw, clean string) (pipelinedb.PipelineErrorRecord, pipelinedb.PipelineCompactErrorRecord) {
	compactClean, filtered := FilterCompactLogText(clean)
	compactRaw, _ := FilterCompactLogText(raw)
	detail := pipelinedb.PipelineErrorRecord{
		RunId: runId, RepoSlug: repo, WorkflowName: workflow,
		StepName: "Failed Step", ErrorText: clean, RawLogs: raw,
	}
	compact := pipelinedb.PipelineCompactErrorRecord{
		RunId: runId, RepoSlug: repo, WorkflowName: workflow,
		StepName: "Failed Step", ErrorText: compactClean, CompactLogs: compactRaw,
		FilteredOkCount: filtered,
	}

	return detail, compact
}

func persistLogToRepoSplitDb(repo string, runId uint64, logContent string) {
	pipeDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return
	}

	defer pipeDb.Close()
	if pipeDb.HasDetailErrorLog(runId) && pipeDb.HasCompactErrorLog(runId) {
		return
	}

	clean := extractCleanErrorLines(logContent)
	workflow := pipeDb.GetRunWorkflowName(runId)
	detail, compact := buildPersistRecords(repo, runId, workflow, logContent, clean)
	_ = pipeDb.RecordDualErrorLog(detail, compact)
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

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func getCachedPipelineLogPathForRepo(repo string, runId uint64) string {
	dir := resolvePipelineDirForRepo(repo)
	repoFile := filepath.Join(dir, fmt.Sprintf("%d.log", runId))
	if isFileExisting(repoFile) {
		return filepath.ToSlash(repoFile)
	}

	legacyFile := filepath.Join(resolvePipelineDir(), fmt.Sprintf("%d.log", runId))
	if isFileExisting(legacyFile) {
		_ = os.Rename(legacyFile, repoFile)

		return filepath.ToSlash(repoFile)
	}

	return filepath.ToSlash(repoFile)
}

func getCachedPipelineJSONPathForRepo(repo string, runId uint64) string {
	dir := resolvePipelineDirForRepo(repo)

	return filepath.ToSlash(filepath.Join(dir, fmt.Sprintf("%d.json", runId)))
}

func resolvePipelineErrorReportPathForRepo(repo string) string {
	dir := resolvePipelineDirForRepo(repo)
	_ = os.MkdirAll(dir, 0755)

	target := filepath.Join(dir, "pipeline_errors.log")
	if isFileExisting(target) {
		return filepath.ToSlash(target)
	}

	legacy := filepath.Join(resolvePipelineDir(), "pipeline_errors.log")
	if isFileExisting(legacy) {
		_ = os.Rename(legacy, target)
	}

	return filepath.ToSlash(target)
}

func writeCombinedErrorReportForRepo(repo string, content string) (string, error) {
	reportPath := resolvePipelineErrorReportPathForRepo(repo)
	if err := os.WriteFile(reportPath, []byte(content), 0644); err != nil {
		return "", err
	}

	_ = writeLastErrorLogForRepo(repo, content)

	return filepath.ToSlash(reportPath), nil
}

func writeLastErrorLogForRepo(repo string, content string) error {
	dir := resolvePipelineDirForRepo(repo)
	_ = os.MkdirAll(dir, 0755)
	lastErrFile := filepath.Join(dir, "last_error.log")

	return os.WriteFile(lastErrFile, []byte(content), 0644)
}

func clearLocalErrorLogs() {
	clearLocalErrorLogsForRepo("")
}

func clearLocalErrorLogsForRepo(repo string) {
	_ = os.Remove(resolvePipelineErrorReportPathForRepo(repo))
	dir := resolvePipelineDirForRepo(repo)
	_ = os.Remove(filepath.Join(dir, "last_error.log"))
	_ = os.Remove(filepath.Join(resolvePipelineDir(), "pipeline_errors.log"))
	_ = os.Remove(filepath.Join(resolvePipelineDir(), "last_error.log"))
}
