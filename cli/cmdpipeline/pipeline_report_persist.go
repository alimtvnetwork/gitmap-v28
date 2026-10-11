package cmdpipeline

import (
	"fmt"
	"os"
	"path/filepath"
)

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
	_ = os.Remove(".gitmap/last_error.log")
	_ = os.Remove(".gitmap/pipeline_errors.log")
	if len(repo) == 0 {
		return
	}
	absPath := findLocalRepoDirForSlug(repo)
	if len(absPath) > 0 {
		_ = os.Remove(filepath.Join(absPath, ".gitmap", "last_error.log"))
		_ = os.Remove(filepath.Join(absPath, ".gitmap", "pipeline_errors.log"))
	}
}
