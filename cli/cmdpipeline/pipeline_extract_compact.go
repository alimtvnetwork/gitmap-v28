package cmdpipeline

import (
	"fmt"
	"strings"
)

func isOkLogLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if isStandardOkLine(trimmed) {
		return true
	}
	if isRustOkLine(trimmed) {
		return true
	}

	return strings.HasPrefix(trimmed, "✔ ok") || strings.HasPrefix(trimmed, "✔ Macro")
}

func isStandardOkLine(trimmed string) bool {
	if trimmed == "PASS" || trimmed == "ok" || strings.HasPrefix(trimmed, "PASS:") {
		return true
	}
	if strings.HasPrefix(trimmed, "ok\t") || strings.HasPrefix(trimmed, "ok ") {
		return true
	}
	if strings.HasPrefix(trimmed, "?\t") || strings.HasPrefix(trimmed, "? ") {
		return true
	}
	if strings.HasPrefix(trimmed, "--- PASS") || strings.HasPrefix(trimmed, "=== RUN") {
		return true
	}

	return false
}

func isRustOkLine(trimmed string) bool {
	if strings.HasPrefix(trimmed, "test result: ok.") {
		return true
	}
	if !strings.HasPrefix(trimmed, "test ") {
		return false
	}
	if strings.HasSuffix(trimmed, " ... ok") || strings.HasSuffix(trimmed, "... ok") {
		return true
	}
	if strings.HasSuffix(trimmed, " ... ignored") || strings.HasSuffix(trimmed, "... ignored") {
		return true
	}

	return false
}

func isKeepLogLine(line string) bool {
	return !isOkLogLine(line) && !isIgnoredLogLine(line)
}

func filterCompactLines(lines []string) []string {
	return ParallelFilterLines(lines, isKeepLogLine)
}

// FilterCompactLogText strips passing ok lines from a text block and returns filtered count.
func FilterCompactLogText(rawText string) (string, int) {
	if len(rawText) == 0 {
		return "", 0
	}

	lines := strings.Split(rawText, "\n")
	filtered := filterCompactLines(lines)
	filteredCount := len(lines) - len(filtered)

	return strings.Join(filtered, "\n"), filteredCount
}

func compactErrorPayload(p *PipelineErrorLogsPayload) {
	for i := range p.FailedRuns {
		compactFailedRunItem(&p.FailedRuns[i])
	}
	for i := range p.SectionFailures {
		compactSectionFailureItem(&p.SectionFailures[i])
	}
	p.CombinedErrors = formatCombinedSectionFailures(p.SectionFailures)
	updateCompactedErrorLogs(p)
}

func compactFailedRunItem(run *FailedRunItem) {
	for j := range run.FailedJobs {
		run.FailedJobs[j].ErrorLines = filterCompactLines(run.FailedJobs[j].ErrorLines)
		run.FailedJobs[j].StackTrace = filterCompactStackTrace(run.FailedJobs[j].StackTrace)
	}
	run.StackTrace = filterCompactStackTrace(run.StackTrace)
}

func compactSectionFailureItem(sec *SectionFailure) {
	sec.ErrorLines = filterCompactLines(sec.ErrorLines)
	sec.StackTrace = filterCompactStackTrace(sec.StackTrace)
}

func updateCompactedErrorLogs(p *PipelineErrorLogsPayload) {
	if len(p.FailedRuns) > 0 {
		p.ErrorLogs = p.CombinedErrors + "\n\n" + formatAllRunsDetailed(p.FailedRuns)

		return
	}
	if len(p.ErrorLogs) > 0 {
		p.ErrorLogs = strings.Join(filterCompactLines(strings.Split(p.ErrorLogs, "\n")), "\n")
	}
}

func isWarningLine(text string) bool {
	trimmed := strings.TrimSpace(text)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "warning:") || strings.HasPrefix(lower, "warning[") || strings.HasPrefix(lower, "warning ") {
		return true
	}
	if strings.Contains(lower, ": warning:") || strings.Contains(lower, ": warning ") {
		return true
	}
	if strings.Contains(lower, "generated ") && strings.Contains(lower, "warnings") {
		return true
	}

	return strings.HasPrefix(lower, "##[warning]")
}

func compactLinkerCommandLine(line string) string {
	if len(line) < 250 {
		return line
	}

	lower := strings.ToLower(line)
	isLinker := strings.Contains(lower, "link.exe") || strings.Contains(lower, "collect2") || strings.Contains(lower, "/libpath:")
	if !isLinker {
		return line
	}

	return formatCompactedLinkerCommand(line)
}

func formatCompactedLinkerCommand(line string) string {
	rlibCount := strings.Count(line, ".rlib")
	libCount := strings.Count(line, ".lib")
	totalLibs := rlibCount + libCount
	if totalLibs < 5 {
		return line
	}

	prefix := extractLinkerPrefix(line)
	outFlag := extractLinkerOutFlag(line)

	return fmt.Sprintf("%s ... [%d library and object files omitted] ... %s", prefix, totalLibs, outFlag)
}

func extractLinkerPrefix(line string) string {
	if idx := strings.Index(line, `"/NOLOGO"`); idx != -1 {
		return line[:idx+len(`"/NOLOGO"`)]
	}
	if idx := strings.Index(line, `link.exe"`); idx != -1 {
		return line[:idx+len(`link.exe"`)]
	}

	return line[:100]
}

func extractLinkerOutFlag(line string) string {
	idx := strings.Index(line, `"/OUT:`)
	if idx == -1 {
		return ""
	}

	endIdx := strings.Index(line[idx+1:], `"`)
	if endIdx != -1 {
		return line[idx : idx+1+endIdx+1]
	}

	return ""
}

func filterOutSummaryLine(lines []string, summary string) []string {
	if len(lines) == 0 {
		return nil
	}

	cleanSummary := strings.ToLower(strings.TrimSpace(summary))
	var out []string
	for _, l := range lines {
		cleanLine := strings.ToLower(strings.TrimSpace(l))
		if cleanLine != cleanSummary && !strings.HasSuffix(cleanLine, cleanSummary) {
			out = append(out, l)
		}
	}

	return out
}
