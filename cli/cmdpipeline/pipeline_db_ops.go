package cmdpipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func extractTargetRepo(args []string) string {
	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		if !strings.HasPrefix(trimmed, "-") && len(trimmed) > 0 {
			return trimmed
		}
	}

	return resolveCurrentRepoSlug()
}

func runPipelineDBClear(args []string) error {
	repo := extractTargetRepo(args)
	if !confirmPipelineDBClear(repo, args) {
		fmt.Println("Clear operation canceled.")

		return nil
	}

	return executePipelineDBClear(repo)
}

func confirmPipelineDBClear(repo string, args []string) bool {
	msg := fmt.Sprintf("Clear all pipeline runs and error logs for %s? [y/N]: ", repo)

	return confirmOrSkip(msg, args)
}

func executePipelineDBClear(repo string) *apperror.AppError {
	dir := resolvePipelineDirForRepo(repo)
	purgedFiles, reclaimedBytes := purgeRepoPipelineFolder(dir)
	clearLocalErrorLogsForRepo(repo)
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err == nil {
		defer db.Close()
		_ = db.Clear()
	}

	printPipelineClearSummary(repo, dir, purgedFiles, reclaimedBytes)

	return nil
}

func purgeRepoPipelineFolder(dir string) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}
	var count int
	var totalBytes int64
	for _, e := range entries {
		if e.IsDir() || !isPurgeablePipelineFile(e.Name()) {
			continue
		}
		c, b := purgeSinglePipelineFile(filepath.Join(dir, e.Name()))
		count += c
		totalBytes += b
	}

	return count, totalBytes
}

func purgeSinglePipelineFile(filePath string) (int, int64) {
	var size int64
	if fi, sErr := os.Stat(filePath); sErr == nil {
		size = fi.Size()
	}
	if rErr := os.Remove(filePath); rErr == nil {
		return 1, size
	}
	return 0, 0
}

func isPurgeablePipelineFile(name string) bool {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".log"):
		return true
	case strings.HasSuffix(lower, ".json"):
		return true
	case strings.HasSuffix(lower, ".db-journal") || strings.HasSuffix(lower, ".db-wal"):
		return true
	default:
		return false
	}
}

func printPipelineClearSummary(repo, dir string, files int, bytes int64) {
	fmt.Printf("%s✓ Pipeline logs and database cleared for %s.%s\n",
		constants.ColorGreen, repo, constants.ColorReset)
	fmt.Printf("  Directory:     %s\n", filepath.ToSlash(dir))
	if files > 0 {
		fmt.Printf("  Files purged:  %d (%s reclaimed)\n", files, formatBytes(bytes))
	} else {
		fmt.Printf("  Files purged:  0 (clean)\n")
	}
	fmt.Printf("  Database:      freshly reset\n")
}

func runPipelineDBReset(args []string) error {
	repo := resolveCurrentRepoSlug()
	if !confirmPipelineDBReset(repo, args) {
		fmt.Println("Reset operation canceled.")

		return nil
	}

	return executePipelineDBReset(repo)
}

func confirmPipelineDBReset(repo string, args []string) bool {
	msg := fmt.Sprintf("Reset and re-create pipeline schema for %s? [y/N]: ", repo)

	return confirmOrSkip(msg, args)
}

func executePipelineDBReset(repo string) *apperror.AppError {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return apperror.WrapSimple(err, "open pipeline split db for reset")
	}
	defer db.Close()

	if resetErr := db.Reset(); resetErr != nil {
		return apperror.WrapSimple(resetErr, "reset pipeline split db")
	}
	fmt.Printf("%s✓ Pipeline split database reset for %s.%s\n", constants.ColorGreen, repo, constants.ColorReset)

	return nil
}

func runPipelineDBOptimize(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return apperror.WrapSimple(err, "open pipeline split db for optimize")
	}
	defer db.Close()

	return executeDBOptimizeAndPrint(db)
}

func executeDBOptimizeAndPrint(db *pipelinedb.PipelineSplitDb) *apperror.AppError {
	reclaimed, err := db.Optimize()
	if err != nil {
		return apperror.WrapSimple(err, "optimize pipeline split db")
	}
	fmt.Printf("%s✓ Pipeline split DB optimized.%s Reclaimed: %s (%s)\n",
		constants.ColorGreen, constants.ColorReset, formatBytes(reclaimed), db.Path)

	return nil
}

func runPipelineDBErrorLogs(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return apperror.WrapSimple(err, "open pipeline split db for error logs")
	}
	defer db.Close()

	return fetchAndRenderDbErrorLogs(db, repo, args)
}

func fetchAndRenderDbErrorLogs(db *pipelinedb.PipelineSplitDb, repo string, args []string) *apperror.AppError {
	logRes := db.QueryRecentErrorLogs(20)
	if logRes.IsFailure() {
		return logRes.AppError()
	}
	if hasArgFlag(args, "--json") {
		return printJSONAppError(logRes.Data)
	}
	if logRes.IsEmpty() {
		fmt.Printf("No error logs recorded in pipeline database for %s.\n", repo)

		return nil
	}
	renderDbErrorLogsList(repo, logRes.Data)

	return nil
}

func printJSONAppError(v any) *apperror.AppError {
	if err := printJSON(v); err != nil {
		return apperror.WrapSimple(err, "render json error logs")
	}

	return nil
}

func renderDbErrorLogsList(repo string, logs []pipelinedb.PipelineErrorRecord) {
	fmt.Printf("\n  %s● Recorded CI/CD Error Logs for %s (%d records):%s\n",
		constants.ColorRed, repo, len(logs), constants.ColorReset)
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
	for _, l := range logs {
		renderSingleDbErrorRecord(l)
	}
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
}

func renderSingleDbErrorRecord(l pipelinedb.PipelineErrorRecord) {
	fmt.Printf("    [%s] Run #%d (%s - %s)\n", l.CreatedAt, l.RunId, l.WorkflowName, l.StepName)
	for _, line := range strings.Split(l.ErrorText, "\n") {
		fmt.Printf("      %s\n", line)
	}
	fmt.Println()
}
