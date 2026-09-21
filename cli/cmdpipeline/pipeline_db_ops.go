package cmdpipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	hasAll := hasArgFlag(args, "--all")
	if hasAll {
		return runPipelineDBClearAll(args)
	}

	repo := extractTargetRepo(args)
	if !confirmPipelineDBClear(repo, args) {
		fmt.Println("Clear operation canceled.")

		return nil
	}

	return executePipelineDBClear(repo)
}

func runPipelineDBClearAll(args []string) error {
	if !confirmPipelineDBClearAll(args) {
		fmt.Println("Clear operation canceled.")

		return nil
	}

	return executePipelineDBClearAll()
}

func confirmPipelineDBClearAll(args []string) bool {
	msg := "Clear all pipeline runs and error logs across ALL repositories? [y/N]: "

	return confirmOrSkip(msg, args)
}

func confirmPipelineDBClear(repo string, args []string) bool {
	msg := fmt.Sprintf("Clear all pipeline runs and error logs for %s? [y/N]: ", repo)

	return confirmOrSkip(msg, args)
}

func executePipelineDBClear(repo string) *apperror.AppError {
	dirs := collectRepoPipelineDirs(repo)
	activeDir := resolvePipelineDirForRepo(repo)
	totalFiles, totalBytes := purgeAllCandidateDirs(dirs, activeDir)
	parentFiles, parentBytes := purgeParentPipelineOrphans(repo)
	totalFiles += parentFiles
	totalBytes += parentBytes

	clearLocalErrorLogsForRepo(repo)
	reinitFreshPipelineDb(repo)
	printPipelineClearSummary(repo, activeDir, totalFiles, totalBytes)

	return nil
}

func executePipelineDBClearAll() *apperror.AppError {
	dirs := []string{
		filepath.Join(store.BinaryDataDir(), "pipeline"),
		resolveGlobalAppDataPipelineDir(),
	}
	totalFiles, totalBytes := purgeAllGlobalDirs(dirs)
	repoRoot := resolveRepoRootDir()
	if len(repoRoot) > 0 && repoRoot != "." {
		localPipe := filepath.Join(repoRoot, ".gitmap", "data", "pipeline")
		f, b := purgeSubtreeAndFiles(localPipe)
		totalFiles += f
		totalBytes += b
	}
	clearLocalErrorLogs()
	reinitFreshPipelineDb(resolveCurrentRepoSlug())
	printClearAllSummary(totalFiles, totalBytes)

	return nil
}

func purgeAllGlobalDirs(dirs []string) (int, int64) {
	seen := make(map[string]bool)
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		clean := filepath.Clean(d)
		if len(d) > 0 && !seen[clean] && isDirExisting(clean) {
			seen[clean] = true
			f, b := purgeSubtreeAndFiles(clean)
			totalFiles += f
			totalBytes += b
		}
	}

	return totalFiles, totalBytes
}

func purgeSubtreeAndFiles(dir string) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}

	var count int
	var bytes int64
	for _, e := range entries {
		c, b := purgeEntryRecursive(dir, e)
		count += c
		bytes += b
	}

	return count, bytes
}

func purgeEntryRecursive(parent string, e os.DirEntry) (int, int64) {
	target := filepath.Join(parent, e.Name())
	if e.IsDir() {
		c, b := purgeRepoPipelineFolder(target)
		_ = os.Remove(target)

		return c, b
	}
	if isPurgeablePipelineFile(e.Name()) {
		return purgeSinglePipelineFile(target)
	}

	return 0, 0
}

func printClearAllSummary(files int, bytes int64) {
	fmt.Printf("%s✓ Pipeline logs and database cleared across ALL repositories.%s\n",
		constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Files purged:  %d (%s reclaimed)\n", files, formatBytes(bytes))
	fmt.Printf("  Database:      freshly reset\n")
}

func purgeAllCandidateDirs(dirs []string, activeDir string) (int, int64) {
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		f, b := purgeRepoPipelineFolder(d)
		totalFiles += f
		totalBytes += b
		if d != activeDir {
			_ = os.Remove(d)
		}
	}

	return totalFiles, totalBytes
}

func reinitFreshPipelineDb(repo string) {
	freshDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err == nil {
		_ = freshDb.Close()
	}
}

func collectRepoPipelineDirs(repo string) []string {
	seen := make(map[string]bool)
	var dirs []string
	slug := pipelinedb.SanitizeRepoSlug(repo)
	underscoreSlug := strings.ReplaceAll(repo, "/", "_")
	candidates := buildDirCandidates(repo, slug, underscoreSlug)

	for _, d := range candidates {
		clean := filepath.Clean(d)
		if isDirExisting(clean) && !seen[clean] {
			seen[clean] = true
			dirs = append(dirs, clean)
		}
	}

	return dirs
}

func resolveGlobalAppDataPipelineDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if len(localAppData) > 0 {
		return filepath.Join(localAppData, "gitmap-cli", "data", "pipeline")
	}

	home, err := os.UserHomeDir()
	if err == nil && len(home) > 0 {
		return filepath.Join(home, "AppData", "Local", "gitmap-cli", "data", "pipeline")
	}

	return ""
}

func buildDirCandidates(repo, slug, underscoreSlug string) []string {
	appData := store.BinaryDataDir()
	candidates := []string{
		pipelinedb.RepoPipelineDir(repo),
		filepath.Join(appData, "pipeline", slug),
		filepath.Join(appData, "pipeline", underscoreSlug),
	}
	candidates = appendGlobalAppDataCandidates(candidates, slug, underscoreSlug)

	return appendRepoRootCandidates(candidates, slug, underscoreSlug)
}

func appendGlobalAppDataCandidates(list []string, slug, underscoreSlug string) []string {
	globalDir := resolveGlobalAppDataPipelineDir()
	if len(globalDir) > 0 {
		list = append(list,
			filepath.Join(globalDir, slug),
			filepath.Join(globalDir, underscoreSlug),
		)
	}

	return list
}

func appendRepoRootCandidates(list []string, slug, underscoreSlug string) []string {
	repoRoot := resolveRepoRootDir()
	if len(repoRoot) > 0 && repoRoot != "." {
		list = append(list,
			filepath.Join(repoRoot, ".gitmap", "data", "pipeline", slug),
			filepath.Join(repoRoot, ".gitmap", "data", "pipeline", underscoreSlug),
		)
	}

	return list
}

func purgeParentPipelineOrphans(repo string) (int, int64) {
	slug := pipelinedb.SanitizeRepoSlug(repo)
	underscoreSlug := strings.ReplaceAll(repo, "/", "_")
	dirs := []string{
		filepath.Join(store.BinaryDataDir(), "pipeline"),
		resolveGlobalAppDataPipelineDir(),
	}

	return purgeOrphanFilesAcrossDirs(dirs, slug, underscoreSlug)
}

func purgeOrphanFilesAcrossDirs(dirs []string, slug, underscoreSlug string) (int, int64) {
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		if len(d) == 0 || !isDirExisting(d) {
			continue
		}
		f, b := purgeOrphansInDir(d, slug, underscoreSlug)
		totalFiles += f
		totalBytes += b
	}

	return totalFiles, totalBytes
}

func purgeOrphansInDir(dir, slug, underscoreSlug string) (int, int64) {
	targets := []string{
		filepath.Join(dir, fmt.Sprintf("pipeline_%s.db", slug)),
		filepath.Join(dir, fmt.Sprintf("%s.db", slug)),
		filepath.Join(dir, fmt.Sprintf("pipeline_%s.db", underscoreSlug)),
		filepath.Join(dir, fmt.Sprintf("%s.db", underscoreSlug)),
	}
	var count int
	var bytes int64
	for _, t := range targets {
		c, b := purgeSinglePipelineFile(t)
		count += c
		bytes += b
	}

	return count, bytes
}

func isDirExisting(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func purgeRepoPipelineFolder(dir string) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}

	return iteratePurgeEntries(dir, entries)
}

func iteratePurgeEntries(dir string, entries []os.DirEntry) (int, int64) {
	var count int
	var totalBytes int64
	for _, e := range entries {
		c, b := purgeEntryIfEligible(dir, e)
		count += c
		totalBytes += b
	}

	return count, totalBytes
}

func purgeEntryIfEligible(dir string, e os.DirEntry) (int, int64) {
	target := filepath.Join(dir, e.Name())
	if e.IsDir() {
		c, b := purgeRepoPipelineFolder(target)
		_ = os.Remove(target)

		return c, b
	}

	return purgeSinglePipelineFile(target)
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
	case strings.HasSuffix(lower, ".db-journal") || strings.HasSuffix(lower, ".db-wal") || strings.HasSuffix(lower, ".db-shm"):
		return true
	case lower == "sql.db" || strings.HasSuffix(lower, ".db"):
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
