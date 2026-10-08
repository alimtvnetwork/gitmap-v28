package cmdpipeline

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderCombinedSectionsTerminal(sections []SectionFailure) {
	if len(sections) == 0 {
		return
	}

	fmt.Printf("  %s● Combined Pipeline Section Failures [%d failed section(s)]:%s\n",
		constants.ColorRed, len(sections), constants.ColorReset)
	for i, sec := range sections {
		renderSingleSectionFailureRow(sec, i+1, len(sections))
	}

	fmt.Println()
}

func renderSingleSectionFailureRow(sec SectionFailure, idx, total int) {
	fmt.Printf("    %s[%d/%d] %s #%d ➔ Job: %s | Step: %s%s\n",
		constants.ColorCyan, idx, total, sec.WorkflowName, sec.RunId, sec.JobName, sec.StepName, constants.ColorReset)
	if len(sec.FailureSummary) > 0 {
		cleaned := cleanDisplayErrorText(sec.FailureSummary, sec.JobName, sec.StepName)
		fmt.Printf("      Error:   %s%s%s\n", constants.ColorRed, cleaned, constants.ColorReset)
	}

	renderSectionSourceDetails(sec)
	renderSectionWarnings(sec.Warnings, activeLogLineLimit)
	renderSectionErrorLinesDedup(sec.ErrorLines, sec.FailureSummary, activeLogLineLimit)
	if len(sec.StackTrace) > 0 {
		renderSectionStackTrace(sec.StackTrace, activeLogLineLimit)
	}

	if len(sec.SavedLogFile) > 0 {
		fmt.Printf("      Log:     %s\n", filepath.ToSlash(FormatRelativeDbPath(sec.SavedLogFile)))
	}
}

func renderSectionSourceDetails(sec SectionFailure) {
	logger := "github-actions"
	if len(sec.WorkflowName) > 0 {
		logger = sec.WorkflowName
	}
	url := fmt.Sprintf("https://github.com/run/%d", sec.RunId)
	logFile := filepath.ToSlash(FormatRelativeDbPath(sec.SavedLogFile))
	if len(logFile) == 0 {
		logFile = "in-memory / active stream"
	}

	fmt.Printf("      Source Details:\n")
	fmt.Printf("        Logger:    %s\n", logger)
	fmt.Printf("        Workflow:  %s\n", sec.WorkflowName)
	fmt.Printf("        Job/Step:  %s / %s\n", sec.JobName, sec.StepName)
	fmt.Printf("        URL:       %s\n", url)
	fmt.Printf("        Log File:  %s\n", logFile)
}

func cleanDisplayErrorText(text, job, step string) string {
	clean := strings.TrimSpace(text)
	prefixBoth := job + "\t" + step + "\t"
	if strings.HasPrefix(clean, prefixBoth) {
		return strings.TrimSpace(strings.TrimPrefix(clean, prefixBoth))
	}

	prefixJob := job + "\t"
	if strings.HasPrefix(clean, prefixJob) {
		return strings.TrimSpace(strings.TrimPrefix(clean, prefixJob))
	}

	return clean
}

func renderSectionStackTrace(stack string, lineLimit ...int) {
	if len(stack) == 0 {
		return
	}

	limit := resolveLogSectionLimit(20, lineLimit...)
	lines := strings.Split(strings.TrimSpace(stack), "\n")
	capped := capErrorLines(lines, limit)

	fmt.Printf("      Stack Trace:\n")
	for _, l := range capped {
		fmt.Printf("        %s%s%s\n", constants.ColorDim, l, constants.ColorReset)
	}

	printRemainingLineCount(len(lines), len(capped))
}

func renderSectionWarnings(warnings []string, lineLimit ...int) {
	if len(warnings) == 0 {
		return
	}

	limit := resolveLogSectionLimit(25, lineLimit...)
	fmt.Printf("      Warnings (%d preceding):\n", len(warnings))
	capped := capErrorLines(warnings, limit)
	for _, w := range capped {
		fmt.Printf("        %s%s%s\n", constants.ColorYellow, w, constants.ColorReset)
	}

	printRemainingLineCount(len(warnings), len(capped))
}

func renderSectionErrorLinesDedup(lines []string, summary string, lineLimit ...int) {
	dedup := filterOutSummaryLine(lines, summary)
	if len(dedup) == 0 {
		return
	}

	limit := resolveLogSectionLimit(6, lineLimit...)
	fmt.Printf("      Details:\n")
	capped := capErrorLines(dedup, limit)
	for _, l := range capped {
		fmt.Printf("        %s%s%s\n", constants.ColorYellow, l, constants.ColorReset)
	}

	printRemainingLineCount(len(dedup), len(capped))
}

func resolveLogSectionLimit(defaultCap int, explicitLimit ...int) int {
	if len(explicitLimit) > 0 && explicitLimit[0] > 0 {
		return explicitLimit[0]
	}
	if activeLogLineLimit > 0 {
		return activeLogLineLimit
	}

	return defaultCap
}

func printRemainingLineCount(total, capped int) {
	if total > capped {
		fmt.Printf("        %s... (%d more lines in log)%s\n", constants.ColorDim, total-capped, constants.ColorReset)
	}
}

func renderFailedRunsBreakdown(failedRuns []FailedRunItem) {
	total := len(failedRuns)
	fmt.Printf("  %s● Detailed Failure Logs Across Workflow Runs [%d run(s)]:%s\n\n",
		constants.ColorRed, total, constants.ColorReset)
	for i, fr := range failedRuns {
		renderFailedRunCard(fr, i+1, total)
	}
}
