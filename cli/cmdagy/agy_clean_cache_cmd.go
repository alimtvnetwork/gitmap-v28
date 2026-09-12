// Package cmdagy — agy_clean_cache_cmd.go defines the CLI command for agy clean-cache.
package cmdagy

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	agyCleanDryRun      bool
	agyCleanYes         bool
	agyCleanForce       bool
	agyCleanNoKill      bool
	agyCleanJSON        bool
	agyCleanIncludeTemp bool
)

var agyCleanCacheCmd = &cobra.Command{
	Use:     "clean-cache",
	Aliases: []string{"cleancache", "clean_cache", "cc"},
	Short:   "Clean Antigravity and browser cache stores and close locking processes",
	Long: `Cleans Chromium, Dawn WebGPU, shader, and scratch caches for Antigravity.
Optionally closes locking processes (antigravity, electron, msedge, msedgewebview2) to ensure
clean unlinking of cached assets.

Protected locations (transcripts, conversation history, user settings, projects)
are strictly preserved.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := CleanCacheOptions{
			DryRun:      agyCleanDryRun,
			Force:       agyCleanForce,
			Yes:         agyCleanYes || agyCleanForce,
			NoKill:      agyCleanNoKill,
			JSON:        agyCleanJSON,
			IncludeTemp: agyCleanIncludeTemp,
		}

		return ExecuteCleanCache(opts)
	},
}

func init() {
	agyCleanCacheCmd.Flags().BoolVarP(&agyCleanDryRun, "dry-run", "d", false, "Preview cache targets and processes without deleting")
	agyCleanCacheCmd.Flags().BoolVarP(&agyCleanYes, "yes", "y", false, "Proceed with cleanup without interactive prompt")
	agyCleanCacheCmd.Flags().BoolVarP(&agyCleanForce, "force", "f", false, "Force terminate processes and proceed without confirmation")
	agyCleanCacheCmd.Flags().BoolVar(&agyCleanNoKill, "no-kill", false, "Skip terminating running processes before cleaning")
	agyCleanCacheCmd.Flags().BoolVar(&agyCleanJSON, "json", false, "Output results in JSON format")
	agyCleanCacheCmd.Flags().BoolVar(&agyCleanIncludeTemp, "include-temp", true, "Include system/user temp directory in cleanup")
}

// ExecuteCleanCache orchestrates cache target discovery, process termination, and cache wiping.
func ExecuteCleanCache(opts CleanCacheOptions) error {
	startTime := time.Now()
	targets := DiscoverCacheTargets(opts.IncludeTemp)
	procs, _ := DiscoverTargetProcesses()

	report := BuildBaseCacheReport(opts, targets, procs)

	if opts.DryRun {
		return handleDryRun(opts, report, targets, procs)
	}

	if !isConfirmed(opts, targets, procs) {
		fmt.Printf("\n%sCleanup canceled. No changes made.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	return executeActualClean(opts, report, targets, procs, startTime)
}

func isConfirmed(opts CleanCacheOptions, targets []AgyCacheTarget, procs []AgyProcessInfo) bool {
	if opts.Yes || opts.JSON {
		return true
	}
	RenderCleanPreview(targets, procs)
	return AskProceedConfirmation()
}

func handleDryRun(opts CleanCacheOptions, report CleanCacheReport, targets []AgyCacheTarget, procs []AgyProcessInfo) error {
	var totalBytes int64
	var totalFiles int

	for _, t := range targets {
		if t.Exists {
			totalBytes += t.SizeBytes
			totalFiles += t.FileCount
		}
	}

	report.BytesFreed = totalBytes
	report.HumanFreed = FormatBytes(totalBytes)
	report.FilesDeleted = totalFiles

	if opts.JSON {
		return RenderCleanJSON(report)
	}

	RenderCleanPreview(targets, procs)
	fmt.Printf("\n%sℹ [dry-run] %d location(s) inspected; %s across %d file(s) would be cleaned.%s\n\n",
		constants.ColorYellow, len(targets), report.HumanFreed, totalFiles, constants.ColorReset)

	return nil
}

func executeActualClean(
	opts CleanCacheOptions,
	report CleanCacheReport,
	targets []AgyCacheTarget,
	procs []AgyProcessInfo,
	startTime time.Time,
) error {
	if isTerminationRequired(opts, procs) {
		logProcessTermination(opts.JSON, len(procs))
		killed, warnings := TerminateProcesses(procs)
		report.ProcessesTerminated = killed
		report.Warnings = append(report.Warnings, warnings...)
	}

	logCacheCleaning(opts.JSON)

	for _, t := range targets {
		if !t.Exists {
			continue
		}
		freed, files, warnings := CleanDirectoryContents(t.Path)
		report.BytesFreed += freed
		report.FilesDeleted += files
		report.Warnings = append(report.Warnings, warnings...)
	}

	report.HumanFreed = FormatBytes(report.BytesFreed)
	report.DurationMs = time.Since(startTime).Milliseconds()

	if opts.JSON {
		return RenderCleanJSON(report)
	}

	RenderCleanSuccess(report)
	return nil
}

func isTerminationRequired(opts CleanCacheOptions, procs []AgyProcessInfo) bool {
	if opts.NoKill || len(procs) == 0 {
		return false
	}
	return true
}

func logProcessTermination(isJSON bool, count int) {
	if isJSON {
		return
	}
	fmt.Printf("\n%sClosing %d process(es)...%s\n", constants.ColorYellow, count, constants.ColorReset)
}

func logCacheCleaning(isJSON bool) {
	if isJSON {
		return
	}
	fmt.Printf("%sCleaning cache directories...%s\n", constants.ColorCyan, constants.ColorReset)
}
