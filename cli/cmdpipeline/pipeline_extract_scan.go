package cmdpipeline

import (
	"bufio"
	"strings"
)

// extractCleanErrorLines filters noisy logs to isolate only failure and error lines.
func extractCleanErrorLines(rawLogs string) string {
	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		return ""
	}

	var parts []string
	for _, j := range jobs {
		lines := append([]string{j.FailureSummary}, j.ErrorLines...)
		cleanLines := deduplicateLines(lines)
		parts = append(parts, strings.Join(cleanLines, "\n"))
	}

	return strings.Join(parts, "\n\n")
}

func deduplicateLines(lines []string) []string {
	var out []string
	seen := make(map[string]bool)
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) > 0 && !seen[trimmed] {
			seen[trimmed] = true
			out = append(out, l)
		}
	}

	return out
}

// ParseFailedLogLines parses raw GitHub Actions log text into structured failed jobs.
func ParseFailedLogLines(rawLogs string) []FailedJobItem {
	if strings.TrimSpace(rawLogs) == "" {
		return nil
	}

	jobMap := make(map[string]*FailedJobItem)
	var order []string
	scanLogLinesIntoMap(rawLogs, jobMap, &order)

	return assembleJobItems(jobMap, order, rawLogs)
}

func scanLogLinesIntoMap(rawLogs string, jobMap map[string]*FailedJobItem, order *[]string) {
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	contextRemaining := 0
	var lastKey string
	warnBuf := make(map[string][]string)
	warnContRemaining := 0
	for scanner.Scan() {
		processLogLine(scanner.Text(), jobMap, order, &contextRemaining, &lastKey, warnBuf, &warnContRemaining)
	}
}

func processLogLine(
	raw string,
	jobMap map[string]*FailedJobItem,
	order *[]string,
	ctxRem *int,
	lastKey *string,
	warnBuf map[string][]string,
	warnContRem *int,
) {
	job, step, text, isError := parseLogLine(raw)
	if text == "" || isIgnoredLogLine(text) {
		*warnContRem = 0

		return
	}

	key := job + "|||" + step
	if isWarningLine(text) {
		bufferWarning(warnBuf, key, text)
		*warnContRem = 10

		return
	}

	if *warnContRem > 0 && isWarningContinuationLine(text) {
		bufferWarning(warnBuf, key, text)
		*warnContRem--

		return
	}
	*warnContRem = 0

	if isError {
		recordErrorLine(jobMap, order, key, job, step, text, ctxRem, lastKey, warnBuf)

		return
	}

	appendContextLine(jobMap, key, text, ctxRem, lastKey)
}

func isWarningContinuationLine(text string) bool {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}
	if strings.HasPrefix(trimmed, "-->") || strings.HasPrefix(trimmed, "::: ") {
		return true
	}
	if strings.Contains(trimmed, " |") || strings.HasPrefix(trimmed, "|") {
		return true
	}
	if strings.HasPrefix(trimmed, "= note:") || strings.HasPrefix(trimmed, "= help:") {
		return true
	}
	if strings.HasPrefix(trimmed, "^") || strings.Contains(trimmed, "^^^") {
		return true
	}

	return false
}

func bufferWarning(warnBuf map[string][]string, key, text string) {
	list := warnBuf[key]
	if len(list) < 30 {
		warnBuf[key] = append(list, text)

		return
	}

	warnBuf[key] = append(list[1:], text)
}

func recordErrorLine(jobMap map[string]*FailedJobItem, order *[]string, key, job, step, text string, ctxRem *int, lastKey *string, warnBuf map[string][]string) {
	item := getOrCreateJobItem(jobMap, order, key, job, step)
	attachBufferedWarnings(item, warnBuf[key])
	compacted := compactLinkerCommandLine(text)
	item.ErrorLines = append(item.ErrorLines, compacted)
	updateJobSummary(item, compacted)
	*ctxRem = resolveContextLimit(text)
	*lastKey = key
}

func attachBufferedWarnings(item *FailedJobItem, warnings []string) {
	if len(item.Warnings) > 0 || len(warnings) == 0 {
		return
	}

	item.Warnings = append(item.Warnings, warnings...)
}

func resolveContextLimit(text string) int {
	if isTerminalStepError(text) {
		return 0
	}

	return 25
}

func isTerminalStepError(text string) bool {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "process completed with exit code") {
		return true
	}
	if strings.Contains(lower, "exit status ") || strings.Contains(lower, "exit code ") {
		return true
	}

	return strings.Contains(lower, "failed (failures=")
}

func appendContextLine(jobMap map[string]*FailedJobItem, key, text string, ctxRem *int, lastKey *string) {
	if *ctxRem <= 0 || *lastKey != key || isOkLogLine(text) {
		return
	}

	item := jobMap[key]
	compacted := compactLinkerCommandLine(text)
	item.ErrorLines = append(item.ErrorLines, "    "+compacted)
	if isWarningLine(text) {
		item.Warnings = append(item.Warnings, text)
	}
	*ctxRem--
}
