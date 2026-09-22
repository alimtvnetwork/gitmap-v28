// Package cmdagy — agy_clean_cache_apply.go applies actual cleanup and conversation pruning.
package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func executeActualClean(
	opts CleanCacheOptions,
	report CleanCacheReport,
	targets []AgyCacheTarget,
	procs []AgyProcessInfo,
	pruned []ConvPruneCandidate,
	startTime time.Time,
) error {
	terminateProcessesIfRequired(opts, procs, &report)
	cleanTargetDirectories(opts, targets, &report)
	pruneConversationsWithStaging(pruned, &report)

	report.HumanFreed = FormatBytes(report.BytesFreed)
	report.DurationMs = time.Since(startTime).Milliseconds()

	if opts.JSON {
		return RenderCleanJSON(report)
	}

	RenderCleanSuccess(report)
	printUndoAdvisoryNotice(report.BackupDir)
	return nil
}

func terminateProcessesIfRequired(opts CleanCacheOptions, procs []AgyProcessInfo, rep *CleanCacheReport) {
	if opts.NoKill || len(procs) == 0 {
		return
	}
	logProcessTermination(opts.JSON, len(procs))
	killed, warnings := TerminateProcesses(procs)
	rep.ProcessesTerminated = killed
	rep.Warnings = append(rep.Warnings, warnings...)
}

func cleanTargetDirectories(opts CleanCacheOptions, targets []AgyCacheTarget, rep *CleanCacheReport) {
	logCacheCleaning(opts.JSON)
	for _, t := range targets {
		if !t.Exists {
			continue
		}
		freed, files, warnings := CleanDirectoryContents(t.Path)
		rep.BytesFreed += freed
		rep.FilesDeleted += files
		rep.Warnings = append(rep.Warnings, warnings...)
	}
}

func pruneConversationsWithStaging(pruned []ConvPruneCandidate, rep *CleanCacheReport) {
	if len(pruned) == 0 {
		return
	}
	backupDir, stageErr := stagePrunedConvs(pruned)
	if stageErr == nil {
		rep.BackupDir = backupDir
	}
	deleted := purgeCandidateConversations(pruned)
	rep.ConvsDeleted = deleted
}

func printUndoAdvisoryNotice(backupDir string) {
	if backupDir == "" {
		return
	}
	fmt.Printf("  %sUndo available: gitmap agy undo (restores from %s)%s\n",
		constants.ColorCyan, backupDir, constants.ColorReset)
	fmt.Printf("  %sWarning: OS temporary directories are ephemeral and may be wiped automatically.%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}
