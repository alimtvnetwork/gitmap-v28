// Package cmdagy — agy_clean_cache_exec.go orchestrates cache cleaning and retention.
package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecuteCleanCache orchestrates cache target discovery, retention pruning, and wiping.
func ExecuteCleanCache(opts CleanCacheOptions) error {
	startTime := time.Now()
	targets := DiscoverCacheTargets(opts.IncludeTemp)
	procs, _ := DiscoverTargetProcesses()
	kept, pruned, _ := evaluateConvRetention(opts.Keep)

	report := BuildBaseCacheReport(opts, targets, procs)
	report.KeepCount = opts.Keep

	if opts.DryRun || opts.Preflight {
		return handlePreflightRun(opts, report, targets, procs, kept, pruned)
	}

	if !isConfirmed(opts, targets, procs, pruned) {
		fmt.Printf("\n%sCleanup canceled. No changes made.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	return executeActualClean(opts, report, targets, procs, pruned, startTime)
}

func handlePreflightRun(opts CleanCacheOptions, rep CleanCacheReport, targets []AgyCacheTarget, procs []AgyProcessInfo, kept, pruned []ConvPruneCandidate) error {
	RenderCleanPreviewWithRetention(targets, procs, kept, pruned)
	if opts.JSON {
		return RenderCleanJSON(rep)
	}
	fmt.Printf("\n%sℹ [preflight] %d cache location(s), %d conversation(s) to prune (keeping %d).%s\n",
		constants.ColorYellow, len(targets), len(pruned), len(kept), constants.ColorReset)
	fmt.Printf("  %sEphemeral temp undo will stage removed items at %s%s\n\n",
		constants.ColorDim, getCacheBackupBaseDir(), constants.ColorReset)
	return nil
}

func isConfirmed(opts CleanCacheOptions, targets []AgyCacheTarget, procs []AgyProcessInfo, pruned []ConvPruneCandidate) bool {
	if opts.Yes || opts.JSON {
		return true
	}
	RenderCleanPreview(targets, procs)
	printPrunePreviewSummary(pruned, opts.Keep)
	return AskProceedConfirmation()
}

func printPrunePreviewSummary(pruned []ConvPruneCandidate, keep int) {
	if len(pruned) > 0 {
		fmt.Printf("  %s• %d conversation(s) exceeding retention (%d) will be archived to temp%s\n",
			constants.ColorYellow, len(pruned), keep, constants.ColorReset)
	}
}

func getCacheBackupBaseDir() string {
	return "os.TempDir/gitmap-agy-cache-backup"
}

// BuildBaseCacheReport constructs an initial report before execution.
func BuildBaseCacheReport(opts CleanCacheOptions, targets []AgyCacheTarget, procs []AgyProcessInfo) CleanCacheReport {
	return CleanCacheReport{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		DryRun:    opts.DryRun,
		Preflight: opts.Preflight,
		Targets:   targets,
		Processes: procs,
	}
}
