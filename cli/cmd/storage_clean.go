package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func runStorageClean(args []string) error {
	opts := parseStorageCleanOptions(args)
	stats, err := executeStorageClean(opts)
	if err != nil {
		return err
	}

	printStorageCleanSummary(stats, opts.IsDryRun)

	return nil
}

func parseStorageCleanOptions(args []string) StorageCleanOptions {
	return StorageCleanOptions{
		IsDryRun:   hasArgFlag(args, "-n") || hasArgFlag(args, "--dry-run"),
		IsVerbose:  hasArgFlag(args, "-v") || hasArgFlag(args, "--verbose"),
		IsForce:    hasArgFlag(args, "-f") || hasArgFlag(args, "--force") || hasArgFlag(args, "-y") || hasArgFlag(args, "--yes"),
		IsVacuumDB: hasArgFlag(args, "-a") || hasArgFlag(args, "--all") || hasArgFlag(args, "--vacuum"),
	}
}

func executeStorageClean(opts StorageCleanOptions) (StorageCleanStats, error) {
	var stats StorageCleanStats
	repoRoot, _ := gitutil.RepoRoot(".")
	if repoRoot == "" {
		return stats, nil
	}

	cleanPipelineLogs(repoRoot, opts, &stats)
	cleanTempDirectory(repoRoot, opts, &stats)
	vacuumIfRequested(repoRoot, opts, &stats)

	return stats, nil
}

func cleanPipelineLogs(repoRoot string, opts StorageCleanOptions, stats *StorageCleanStats) {
	pipelineDir := filepath.Join(repoRoot, ".gitmap", "pipeline")
	entries, err := os.ReadDir(pipelineDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		processPipelineLogEntry(pipelineDir, e, opts, stats)
	}
}

func processPipelineLogEntry(dir string, e os.DirEntry, opts StorageCleanOptions, stats *StorageCleanStats) {
	if e.IsDir() {
		return
	}
	name := e.Name()
	isTarget := strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".json")
	if isTarget {
		cleanSingleFile(filepath.Join(dir, name), opts, stats)
		stats.DeletedLogsCount++
	}
}

func cleanTempDirectory(repoRoot string, opts StorageCleanOptions, stats *StorageCleanStats) {
	tempDir := filepath.Join(repoRoot, ".lovable", "temp")
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		processTempEntry(tempDir, e, opts, stats)
	}
}

func processTempEntry(dir string, e os.DirEntry, opts StorageCleanOptions, stats *StorageCleanStats) {
	if e.IsDir() {
		return
	}
	cleanSingleFile(filepath.Join(dir, e.Name()), opts, stats)
	stats.DeletedTempCount++
}

func cleanSingleFile(path string, opts StorageCleanOptions, stats *StorageCleanStats) {
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

func vacuumIfRequested(repoRoot string, opts StorageCleanOptions, stats *StorageCleanStats) {
	if !opts.IsVacuumDB || opts.IsDryRun {
		return
	}

	repoSlug := filepath.Base(repoRoot)
	db, err := pipelinedb.OpenPipelineSplitDb(repoSlug)
	if err != nil {
		return
	}
	defer db.Close()
	freed, _ := db.Vacuum()
	stats.VacuumFreedBytes = freed
	stats.ReclaimedBytes += freed
}

func printStorageCleanSummary(stats StorageCleanStats, isDryRun bool) {
	prefix := "✓ Cleaned"
	if isDryRun {
		prefix = "✔ [Dry-run] Would clean"
	}
	totalFiles := stats.DeletedLogsCount + stats.DeletedTempCount
	bytesStr := cmddb.FormatBytes(stats.ReclaimedBytes)
	fmt.Printf("\n  %s %d log file(s), %d temp file(s) (%s reclaimable)\n\n",
		prefix, stats.DeletedLogsCount, stats.DeletedTempCount, bytesStr)
	if stats.VacuumFreedBytes > 0 {
		fmt.Printf("  Vacuum reclaimed: %s\n\n", cmddb.FormatBytes(stats.VacuumFreedBytes))
	}
	_ = totalFiles
}
