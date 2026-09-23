package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runStorageResetErrors(args []string) error {
	opts := parseStorageResetOptions(args)
	stats := executeStorageReset(opts)
	printStorageResetSummary(stats, opts.IsDryRun)

	return nil
}

func parseStorageResetOptions(args []string) StorageResetOptions {
	return StorageResetOptions{
		IsDryRun:  hasArgFlag(args, "-n") || hasArgFlag(args, "--dry-run"),
		IsVerbose: hasArgFlag(args, "-v") || hasArgFlag(args, "--verbose"),
		IsForce:   hasArgFlag(args, "-f") || hasArgFlag(args, "--force") || hasArgFlag(args, "-y") || hasArgFlag(args, "--yes"),
	}
}

func executeStorageReset(opts StorageResetOptions) StorageResetStats {
	var stats StorageResetStats
	repoRoot := resolveCurrentRepoRoot()
	purgePipelineLogs(repoRoot, opts, &stats)
	purgeScanReports(repoRoot, opts, &stats)
	purgeDatabaseErrors(repoRoot, opts, &stats)
	purgeAGYBackups(opts, &stats)

	return stats
}

func purgeAGYBackups(opts StorageResetOptions, stats *StorageResetStats) {
	baseDir := filepath.Join(store.BinaryDataDir(), "agy")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return
	}
	_ = filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			stats.ReclaimedBytes += info.Size()
		}
		return nil
	})
	if !opts.IsDryRun {
		_ = os.RemoveAll(baseDir)
	}
}

func resolveCurrentRepoRoot() string {
	root, err := gitutil.RepoRoot(".")
	if err != nil || len(root) == 0 {
		return "."
	}

	return root
}

func purgePipelineLogs(repoRoot string, opts StorageResetOptions, stats *StorageResetStats) {
	dirs := candidatePipelineDirs(repoRoot)
	for _, dir := range dirs {
		purgeDirLogFiles(dir, opts, stats)
	}
	purgeLastErrorLog(repoRoot, opts, stats)
}

func candidatePipelineDirs(repoRoot string) []string {
	list := []string{filepath.Join(repoRoot, ".gitmap", "pipeline")}
	home, err := os.UserHomeDir()
	if err == nil {
		list = append(list, filepath.Join(home, ".gitmap", "pipeline"))
	}

	return list
}

func purgeLastErrorLog(repoRoot string, opts StorageResetOptions, stats *StorageResetStats) {
	logFile := filepath.Join(repoRoot, ".gitmap", "last_error.log")
	purgeSingleErrorFile(logFile, opts, stats)
}

func purgeDirLogFiles(dir string, opts StorageResetOptions, stats *StorageResetStats) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		processPipelineDirEntry(dir, e, opts, stats)
	}
}

func processPipelineDirEntry(dir string, e os.DirEntry, opts StorageResetOptions, stats *StorageResetStats) {
	if e.IsDir() {
		return
	}
	if isErrorLogOrJsonFile(e.Name()) {
		purgeSingleErrorFile(filepath.Join(dir, e.Name()), opts, stats)
		stats.ClearedPipelineLogs++
	}
}

func isErrorLogOrJsonFile(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasSuffix(lower, ".log") || strings.HasSuffix(lower, ".json")
}

func purgeScanReports(repoRoot string, opts StorageResetOptions, stats *StorageResetStats) {
	dirs := candidateReportDirs(repoRoot)
	for _, dir := range dirs {
		purgeReportsInDir(dir, opts, stats)
	}
}

func candidateReportDirs(repoRoot string) []string {
	list := []string{filepath.Join(repoRoot, ".gitmap", "reports")}
	binaryReports := filepath.Join(store.BinaryDataDir(), "reports")
	list = append(list, binaryReports)

	return list
}

func purgeReportsInDir(dir string, opts StorageResetOptions, stats *StorageResetStats) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		processReportDirEntry(dir, e, opts, stats)
	}
}

func processReportDirEntry(dir string, e os.DirEntry, opts StorageResetOptions, stats *StorageResetStats) {
	if e.IsDir() {
		return
	}
	if isScanReportFile(e.Name()) {
		purgeSingleErrorFile(filepath.Join(dir, e.Name()), opts, stats)
		stats.ClearedErrorReports++
	}
}

func isScanReportFile(name string) bool {
	lower := strings.ToLower(name)

	return strings.HasPrefix(lower, "errors-") || strings.HasSuffix(lower, ".json")
}

func purgeSingleErrorFile(path string, opts StorageResetOptions, stats *StorageResetStats) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	stats.ReclaimedBytes += fi.Size()
	if opts.IsVerbose {
		fmt.Printf("  • %s (%s)\n", path, cmddb.FormatBytes(fi.Size()))
	}
	if !opts.IsDryRun {
		_ = os.Remove(path)
	}
}

func purgeDatabaseErrors(repoRoot string, opts StorageResetOptions, stats *StorageResetStats) {
	dbPaths := candidatePipelineDbPaths(repoRoot)
	for _, dbPath := range dbPaths {
		purgeOneDatabaseErrors(dbPath, opts, stats)
	}
}

func candidatePipelineDbPaths(repoRoot string) []string {
	var paths []string
	dirs := candidatePipelineDirs(repoRoot)
	for _, dir := range dirs {
		paths = append(paths, findDbFilesInDir(dir)...)
	}

	return paths
}

func findDbFilesInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var matched []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".db") {
			matched = append(matched, filepath.Join(dir, e.Name()))
		}
	}

	return matched
}

func purgeOneDatabaseErrors(dbPath string, opts StorageResetOptions, stats *StorageResetStats) {
	if opts.IsDryRun {
		stats.ClearedDbRecords += countDatabaseErrors(dbPath)
		return
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return
	}
	defer db.Close()

	deleted := executeErrorTableTruncates(db)
	stats.ClearedDbRecords += deleted
	vacuumFreed := runDbVacuum(db, dbPath)
	stats.ReclaimedBytes += vacuumFreed
}

func countDatabaseErrors(dbPath string) int {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0
	}
	defer db.Close()

	total := 0
	tables := []string{"PipelineErrorRecord", "PipelineDetailErrorLog", "PipelineCompactErrorLog", "PipelineErrorLog"}
	for _, tbl := range tables {
		total += countSingleTable(db, tbl)
	}

	return total
}

func countSingleTable(db *sql.DB, tableName string) int {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s;", tableName)
	if err := db.QueryRow(query).Scan(&count); err != nil {
		return 0
	}

	return count
}

func executeErrorTableTruncates(db *sql.DB) int {
	tables := []string{"PipelineErrorRecord", "PipelineDetailErrorLog", "PipelineCompactErrorLog", "PipelineErrorLog"}
	deletedTotal := 0
	for _, tbl := range tables {
		deletedTotal += truncateOneErrorTable(db, tbl)
	}
	resetErrorSequence(db)

	return deletedTotal
}

func truncateOneErrorTable(db *sql.DB, tableName string) int {
	count := countSingleTable(db, tableName)
	query := fmt.Sprintf("DELETE FROM %s;", tableName)
	if _, err := db.Exec(query); err != nil {
		return 0
	}

	return count
}

func resetErrorSequence(db *sql.DB) {
	q := "DELETE FROM sqlite_sequence WHERE name IN ('PipelineCompactErrorLog', 'PipelineDetailErrorLog', 'PipelineErrorLog', 'PipelineErrorRecord');"
	if _, err := db.Exec(q); err != nil {
		return
	}
}

func runDbVacuum(db *sql.DB, path string) int64 {
	beforeSize := fileSizeBytes(path)
	if _, err := db.Exec("VACUUM;"); err != nil {
		return 0
	}
	afterSize := fileSizeBytes(path)
	if beforeSize > afterSize {
		return beforeSize - afterSize
	}

	return 0
}

func fileSizeBytes(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func printStorageResetSummary(stats StorageResetStats, isDryRun bool) {
	title := "✓ Reset error storage"
	if isDryRun {
		title = "✔ [Dry-run] Would reset error storage"
	}
	fmt.Printf("\n  %s:\n", title)
	fmt.Printf("    • Pipeline error logs:  %d file(s)\n", stats.ClearedPipelineLogs)
	fmt.Printf("    • Scan error reports:   %d file(s)\n", stats.ClearedErrorReports)
	fmt.Printf("    • Database error rows:  %d record(s)\n", stats.ClearedDbRecords)
	if stats.ReclaimedBytes > 0 {
		fmt.Printf("    • Space reclaimed:      %s\n", cmddb.FormatBytes(stats.ReclaimedBytes))
	}
	fmt.Println()
}
