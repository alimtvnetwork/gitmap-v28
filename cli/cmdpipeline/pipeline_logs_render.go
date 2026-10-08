package cmdpipeline

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderErrorLogsTerminal(p PipelineErrorLogsPayload) {
	if p.IsRunning {
		renderActiveRunningBanner(p)
	}

	if p.IsRunning && len(p.FailedRuns) == 0 {
		renderRunningSummaryTable(p)
		copyReportToClipboard(buildClipboardCleanReport(p), false)

		return
	}

	renderAndCopyTerminalReport(p)
}

func renderRunningSummaryTable(p PipelineErrorLogsPayload) {
	runs := p.Runs
	if len(runs) == 0 {
		runs = resolveCachedRunsOrFetch(p.Repo)
	}
	RenderHistorySummaryTable(runs)
}

func renderAndCopyTerminalReport(p PipelineErrorLogsPayload) {
	hasFailure := isFailingConclusion(p.Conclusion) || len(p.FailedRuns) > 0
	if hasFailure {
		renderFailureTerminal(p)
		copyReportToClipboard(buildClipboardErrorReport(p), true)

		return
	}

	renderCleanSuccessTerminal(p)
	copyReportToClipboard(buildClipboardCleanReport(p), false)
}

func renderCleanSuccessTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Status: CLEAN (No errors found)%s\n",
		constants.ColorGreen, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Printf("\n  %s✓ All recent pipeline runs for %s on branch %s are PASSING (100%% green).%s\n\n",
		constants.ColorGreen, p.Repo, p.LatestBranch, constants.ColorReset)
	renderCleanSuccessDbAndHistory(p)
	printRerunETA(p.RerunEtaSeconds)
}

func renderPipelineMetaBlock(p PipelineErrorLogsPayload) {
	if p.IsFromCache {
		renderCacheBannerLine(p)
	}
	renderPipelineRepoAndBranch(p)
	renderPipelineMetaVersionAndPR(p)
}

func renderPipelineRepoAndBranch(p PipelineErrorLogsPayload) {
	fmt.Printf("    %-18s %s\n", "Repo:", p.Repo)
	if len(p.RepoUrl) > 0 {
		fmt.Printf("    %-18s %s\n", "Repo URL:", p.RepoUrl)
	}
	if len(p.LatestBranch) > 0 {
		fmt.Printf("    %-18s %s\n", "Branch:", p.LatestBranch)
	}
	if len(p.LastHash) > 0 {
		fmt.Printf("    %-18s %s\n", "Last Commit:", p.LastHash)
	}
}

func renderCacheBannerLine(p PipelineErrorLogsPayload) {
	source := p.CacheSource
	if len(source) == 0 {
		source = "sqlite"
	}
	fmt.Printf("    %-18s %s⚡ Served from local %s DB cache (commit %s)%s\n",
		"Cache:", constants.ColorCyan, strings.ToUpper(source), p.LastHash, constants.ColorReset)
}

func renderPipelineMetaVersionAndPR(p PipelineErrorLogsPayload) {
	if len(p.LastReleaseVersion) > 0 {
		fmt.Printf("    %-18s %s\n", "Last Release:", p.LastReleaseVersion)
	}
	fmt.Printf("    %-18s %d\n", "Open PRs:", p.OpenPRsCount)
	if len(p.Status) > 0 && len(p.Conclusion) > 0 {
		fmt.Printf("    %-18s %s (conclusion: %s)\n", "Run Status:", p.Status, p.Conclusion)
	}
}

func renderCleanSuccessDbAndHistory(p PipelineErrorLogsPayload) {
	if len(p.DbPath) > 0 {
		fmt.Printf("  • Pipeline DB:     %s\n", filepath.ToSlash(FormatRelativeDbPath(p.DbPath)))
		fmt.Printf("  • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
		fmt.Printf("  • Cleanup:         gitmap pipeline clear -y\n")
	}

	runs := p.Runs
	if len(runs) == 0 {
		runs = resolveCachedRunsOrFetch(p.Repo)
	}
	RenderHistorySummaryTable(runs)
}

func renderActiveRunningBanner(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Active Pipeline is RUNNING%s: [%s #%d] (ETA: %s)\n",
		constants.ColorYellow, constants.ColorReset, p.ActiveRunName, p.ActiveRunId, formatEtaDisplay(p.EtaSeconds))
	if len(p.ActiveRunUrl) > 0 {
		fmt.Printf("    URL: %s\n\n", p.ActiveRunUrl)
	}
}

func renderFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Failure Detected%s\n",
		constants.ColorRed, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Println()
	renderFailureTreeTerminal(p)
	if len(p.FailedRuns) == 0 {
		renderSingleFailureTerminal(p)

		return
	}

	renderFailureSectionsAndETA(p)
}

func renderFailureSectionsAndETA(p PipelineErrorLogsPayload) {
	renderCombinedSectionsTerminal(p.SectionFailures)
	renderFailedRunsBreakdown(p.FailedRuns)
	renderSavedLocationsTerminal(p)
	runs := p.Runs
	if len(runs) == 0 {
		runs = resolveCachedRunsOrFetch(p.Repo)
	}
	RenderHistorySummaryTable(runs)
	printRerunETA(p.RerunEtaSeconds)
}
