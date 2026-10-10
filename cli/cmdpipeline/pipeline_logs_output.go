package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"path/filepath"
	"strings"
)

func writeOrRenderErrorLogs(params ErrorLogOutputParams) error {
	contentToWrite, err := formatErrorLogContent(params)
	if err != nil {
		return err
	}
	_ = persistAutoErrorReport(params)

	if len(params.TempFile) > 0 || len(params.FilePath) > 0 {
		return writeErrorLogsToDisk(params, contentToWrite)
	}

	return dispatchErrorLogPresentation(params, contentToWrite)
}

func dispatchErrorLogPresentation(params ErrorLogOutputParams, content string) error {
	if params.IsJSON {
		return outputJSONErrorLogs(content)
	}
	if params.HasSuppressOutputLog {
		printSuppressedStagingNotice(params.Payload.SavedReportFile)

		return nil
	}
	renderErrorLogsTerminal(params.Payload)

	hasFailure := params.Payload.Conclusion == "failure" || len(params.Payload.FailedRuns) > 0
	if hasFailure && params.WantFix && PipelineAgyFixRunner != nil {
		fmt.Printf("\n  🚀 Automatically dispatching CI/CD fix to Antigravity IDE...\n")
		return PipelineAgyFixRunner([]string{params.Payload.Repo, "--force"})
	}

	return nil
}

func outputJSONErrorLogs(content string) error {
	fmt.Println(content)
	if isClipboardWriteAllowed(false) {
		_ = writeClipboard(content)
	}

	return nil
}

func printSuppressedStagingNotice(reportFile string) {
	if len(reportFile) > 0 {
		fmt.Printf("  ✓ Error logs staged to %s (suppressed terminal output via --no-output-log)\n", reportFile)

		return
	}

	fmt.Println("  ✓ Error logs staged to filesystem (suppressed terminal output via --no-output-log)")
}

func writeErrorLogsToDisk(params ErrorLogOutputParams, content string) error {
	if isClipboardWriteAllowed(false) {
		_ = writeClipboard(content)
	}
	if len(params.TempFile) > 0 {
		targetPath := filepath.Join(resolveTempDir(), params.TempFile)

		return writeContentToFile(targetPath, content)
	}

	return writeContentToFile(params.FilePath, content)
}

func persistAutoErrorReport(params ErrorLogOutputParams) error {
	if !isFailingConclusion(params.Payload.Conclusion) && len(params.Payload.FailedRuns) == 0 {
		clearLocalErrorLogs()

		return nil
	}

	return saveActiveErrorReport(params.Payload)
}

func saveActiveErrorReport(p PipelineErrorLogsPayload) error {
	reportContent := p.ErrorLogs
	if len(reportContent) == 0 {
		reportContent = p.CombinedErrors
	}

	if len(reportContent) == 0 {
		return nil
	}

	_, err := writeCombinedErrorReportForRepo(p.Repo, reportContent)

	return err
}

func formatErrorLogContent(params ErrorLogOutputParams) (string, error) {
	if !params.IsJSON {
		return params.Payload.ErrorLogs, nil
	}

	b, err := json.MarshalIndent(params.Payload, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

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
