package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	pipelineExecutionStartTime time.Time
	activeCommandTag           string
)

func markPipelineExecutionStart() {
	if pipelineExecutionStartTime.IsZero() {
		pipelineExecutionStartTime = time.Now()
	}
}

func SetPipelineExecutionStartTime(t time.Time) {
	pipelineExecutionStartTime = t
}

func SetActiveCommandTagForTest(tag string) {
	activeCommandTag = tag
}

func getPipelineExecutionElapsed() time.Duration {
	if pipelineExecutionStartTime.IsZero() {
		return 15 * time.Millisecond
	}
	return time.Since(pipelineExecutionStartTime)
}

func attachExecutionDuration(payload *PipelineErrorLogsPayload) {
	if payload.ElapsedSeconds > 0 && len(payload.ElapsedDuration) > 0 {
		return
	}
	elapsed := getPipelineExecutionElapsed()
	payload.ElapsedSeconds = elapsed.Seconds()
	payload.ElapsedDuration = formatExecutionDuration(elapsed)
}

func formatExecutionDuration(d time.Duration) string {
	if d < time.Second {
		ms := d.Milliseconds()
		if ms <= 0 {
			ms = 1
		}
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func matchArgCommandTag(arg string) string {
	lower := strings.ToLower(arg)
	if lower == "pe" || lower == "pipeline-errors" || lower == "pipeline_errors" || lower == "ee" {
		return "[gitmap pe]"
	}
	if lower == "te" {
		return "[gitmap te]"
	}
	return ""
}

func resolveCommandTag() string {
	if len(activeCommandTag) > 0 {
		return activeCommandTag
	}
	for _, arg := range os.Args {
		if tag := matchArgCommandTag(arg); tag != "" {
			return tag
		}
	}
	return "[gitmap te]"
}

func formatTimingBanner(cmdTag string, elapsed time.Duration, isFromCache bool) string {
	dur := formatExecutionDuration(elapsed)
	if isFromCache {
		return fmt.Sprintf("⚡ %s Request completed in %s (from cache)", cmdTag, dur)
	}
	return fmt.Sprintf("⏱️ %s Request completed in %s", cmdTag, dur)
}

func FormatTimingBanner(p PipelineErrorLogsPayload) string {
	cmdTag := resolveCommandTag()
	elapsed := getPipelineExecutionElapsed()
	if p.ElapsedSeconds > 0 {
		elapsed = time.Duration(p.ElapsedSeconds * float64(time.Second))
	}
	return formatTimingBanner(cmdTag, elapsed, p.IsFromCache)
}

func printTimingBanner(p PipelineErrorLogsPayload) {
	fmt.Printf("\n  %s\n", FormatTimingBanner(p))
}

func writeOrRenderErrorLogs(params ErrorLogOutputParams) error {
	attachExecutionDuration(&params.Payload)
	contentToWrite, err := formatErrorLogContent(params)
	if err != nil {
		return err
	}

	if err := persistAutoErrorReport(params); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not persist error report: %v\n", err)
	}

	if len(params.TempFile) > 0 || len(params.FilePath) > 0 {
		return writeErrorLogsToDisk(params, contentToWrite)
	}

	return dispatchErrorLogPresentation(params, contentToWrite)
}

func dispatchErrorLogPresentation(params ErrorLogOutputParams, content string) error {
	if params.IsJSON {
		return outputJSONErrorLogs(content)
	}

	defer printTimingBanner(params.Payload)
	if params.HasSuppressOutputLog {
		printSuppressedStagingNotice(params.Payload.SavedReportFile)

		return nil
	}

	renderErrorLogsTerminal(params.Payload)

	return maybeDispatchAgyFix(params)
}

func maybeDispatchAgyFix(params ErrorLogOutputParams) error {
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
		if err := writeClipboard(content); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: could not copy to clipboard: %v\n", err)
		}
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
		if err := writeClipboard(content); err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: could not copy to clipboard: %v\n", err)
		}
	}

	defer printTimingBanner(params.Payload)
	if len(params.TempFile) > 0 {
		return writeContentToFile(filepath.Join(resolveTempDir(), params.TempFile), content)
	}

	return writeContentToFile(params.FilePath, content)
}

func persistAutoErrorReport(params ErrorLogOutputParams) error {
	if !isFailingConclusion(params.Payload.Conclusion) && len(params.Payload.FailedRuns) == 0 {
		clearLocalErrorLogs()
		return nil
	}
	report := params.Payload.ErrorLogs
	if len(report) == 0 {
		report = params.Payload.CombinedErrors
	}
	if len(report) == 0 {
		return nil
	}
	_, err := writeCombinedErrorReportForRepo(params.Payload.Repo, report)
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
		renderRunsSummaryTable(p)
		copyReportToClipboard(buildClipboardCleanReport(p), false)
		return
	}
	renderAndCopyTerminalReport(p)
}

func renderRunsSummaryTable(p PipelineErrorLogsPayload) {
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
	if len(p.DbPath) > 0 {
		fmt.Printf("  • Pipeline DB:     %s\n", filepath.ToSlash(FormatRelativeDbPath(p.DbPath)))
		fmt.Printf("  • DB Size:         %s\n", ResolveDbFileSize(p.DbPath))
		fmt.Printf("  • Cleanup:         gitmap pipeline clear -y\n")
	}
	renderRunsSummaryTable(p)
	printRerunETA(p.RerunEtaSeconds)
}

func printMetaLine(label, val string) {
	if len(val) > 0 {
		fmt.Printf("    %-18s %s\n", label, val)
	}
}

func renderPipelineMetaBlock(p PipelineErrorLogsPayload) {
	if p.IsFromCache {
		renderCacheBannerLine(p)
	}
	printMetaLine("Repo:", p.Repo)
	printMetaLine("Repo URL:", p.RepoUrl)
	printMetaLine("Branch:", p.LatestBranch)
	printMetaLine("Last Commit:", p.LastHash)
	printMetaLine("Last Release:", p.LastReleaseVersion)
	fmt.Printf("    %-18s %d\n", "Open PRs:", p.OpenPRsCount)
	if len(p.Status) > 0 && len(p.Conclusion) > 0 {
		fmt.Printf("    %-18s %s (conclusion: %s)\n", "Run Status:", p.Status, p.Conclusion)
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

func renderActiveRunningBanner(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Active Pipeline is RUNNING%s: [%s #%d] (ETA: %s)\n",
		constants.ColorYellow, constants.ColorReset, p.ActiveRunName, p.ActiveRunId, formatEtaDisplay(p.EtaSeconds))
	if len(p.ActiveRunUrl) > 0 {
		fmt.Printf("    URL: %s\n\n", p.ActiveRunUrl)
	}
}

func renderFailureTerminal(p PipelineErrorLogsPayload) {
	fmt.Printf("  %s● Pipeline Failure Detected%s\n", constants.ColorRed, constants.ColorReset)
	renderPipelineMetaBlock(p)
	fmt.Println()
	renderFailureTreeTerminal(p)
	if len(p.FailedRuns) == 0 {
		renderSingleFailureTerminal(p)
		return
	}
	renderCombinedSectionsTerminal(p.SectionFailures)
	renderFailedRunsBreakdown(p.FailedRuns)
	renderSavedLocationsTerminal(p)
	renderRunsSummaryTable(p)
	printRerunETA(p.RerunEtaSeconds)
}
