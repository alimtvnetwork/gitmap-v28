package cmdpipeline

import (
	"bufio"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

func extractStackTraceFromLog(rawLogs string) string {
	startIdx := findStackTraceStart(rawLogs)
	if startIdx < 0 {
		return ""
	}

	return extractBoundedStackLines(rawLogs[startIdx:])
}

func extractJobStackTrace(rawLogs, jobName string) string {
	if len(rawLogs) == 0 || len(jobName) == 0 {
		return ""
	}

	jobLogs := extractLinesForJob(rawLogs, jobName)
	if len(jobLogs) == 0 {
		return ""
	}

	return extractStackTraceFromLog(jobLogs)
}

func extractLinesForJob(rawLogs, jobName string) string {
	var matched []string
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	for scanner.Scan() {
		raw := scanner.Text()
		job, _, _, _ := parseLogLine(raw)
		if isMatchingJobName(job, jobName) && len(strings.TrimSpace(raw)) > 0 {
			matched = append(matched, raw)
		}
	}

	return strings.Join(matched, "\n")
}

var (
	pyFrameRegex     = lazyregex.New(`File\s+"[^"]+",\s+line\s+\d+`)
	pyExceptionRegex = lazyregex.New(`^(?:[A-Za-z_][A-Za-z0-9_.]*\.)?[A-Za-z_][A-Za-z0-9_]*(?:Error|Exception|Exit|Interrupt|Failure)(?::\s*.*)?$`)
)

func findStackTraceStart(raw string) int {
	if tbIdx := strings.Index(raw, "Traceback (most recent call last):"); tbIdx != -1 {
		return resolveTracebackStartIndex(raw, tbIdx)
	}
	if gIdx := strings.Index(raw, "goroutine "); gIdx != -1 {
		return gIdx
	}
	if pIdx := strings.Index(raw, "panic:"); pIdx != -1 {
		return pIdx
	}
	if sIdx := strings.Index(raw, "Stack Trace:"); sIdx != -1 {
		return sIdx + len("Stack Trace:")
	}

	return -1
}

func resolveTracebackStartIndex(raw string, tbIdx int) int {
	headerIdx := findPrecedingTestHeaderIndex(raw, tbIdx)
	if headerIdx >= 0 {
		return headerIdx
	}

	return tbIdx
}

func findPrecedingTestHeaderIndex(raw string, tbIdx int) int {
	searchStart := 0
	if tbIdx > 500 {
		searchStart = tbIdx - 500
	}
	preceding := raw[searchStart:tbIdx]
	bestOffset := findLatestHeaderOffset(preceding)
	if bestOffset < 0 {
		return -1
	}

	absIdx := searchStart + bestOffset
	lineStart := strings.LastIndex(raw[:absIdx], "\n")
	if lineStart != -1 {
		return lineStart + 1
	}

	return 0
}

func findLatestHeaderOffset(preceding string) int {
	markers := []string{"FAIL: ", "ERROR: ", "--- FAIL: "}
	bestOffset := -1
	for _, m := range markers {
		idx := strings.LastIndex(preceding, m)
		if idx > bestOffset {
			bestOffset = idx
		}
	}

	return bestOffset
}

func extractBoundedStackLines(sub string) string {
	scanner := bufio.NewScanner(strings.NewReader(sub))
	var frames []string
	hasHeader := false
	expectPyCode := false
	for scanner.Scan() {
		if stop := processStackScanLine(scanner.Text(), &frames, &hasHeader, &expectPyCode); stop {
			break
		}
	}

	return strings.Join(frames, "\n")
}

func processStackScanLine(raw string, frames *[]string, hasHeader *bool, expectPyCode *bool) bool {
	line := sanitizeStackLine(raw)
	if *hasHeader && isStackStopLine(line, len(*frames), *expectPyCode) {
		return true
	}
	if !*hasHeader && isStackStartLine(line) {
		*hasHeader = true
		*frames = append(*frames, line)
		return false
	}
	if !*hasHeader || len(strings.TrimSpace(line)) == 0 {
		return false
	}
	if *expectPyCode {
		*expectPyCode = false
		*frames = append(*frames, line)
		return len(*frames) >= 25
	}
	if isPythonStackFrame(line) {
		*expectPyCode = true
	}
	*frames = append(*frames, line)

	return len(*frames) >= 25
}

func isStackStopLine(line string, frameCount int, expectPyCode bool) bool {
	if isStackTerminator(line) {
		return true
	}
	if len(strings.TrimSpace(line)) == 0 {
		return false
	}
	if expectPyCode {
		return false
	}

	return !isStackFrame(line)
}

func sanitizeStackLine(raw string) string {
	clean := ansiRegex.ReplaceAllString(raw, "")
	parts := strings.Split(clean, "\t")
	var text string
	if len(parts) >= 3 {
		text = stripTimestamp(strings.Join(parts[2:], "\t"))
	} else if len(parts) == 2 {
		text = parts[1]
	} else {
		text = clean
	}
	text = strings.TrimPrefix(text, "##[error]")
	text = strings.TrimPrefix(text, "##[group]")
	text = strings.TrimPrefix(text, "##[endgroup]")
	text = cleanAnnotationError(text)

	return strings.TrimRight(text, "\r\n")
}

func isStackStartLine(line string) bool {
	if strings.Contains(line, "Traceback (most recent call last):") {
		return true
	}
	if isTestHeaderStartLine(line) {
		return true
	}

	return strings.HasPrefix(line, "goroutine ") || strings.HasPrefix(line, "panic:") || strings.Contains(line, "panic:")
}

func isTestHeaderStartLine(line string) bool {
	if strings.HasPrefix(line, "FAIL: ") || strings.HasPrefix(line, "ERROR: ") {
		return true
	}

	return strings.HasPrefix(line, "--- FAIL:")
}

func isStackFrame(line string) bool {
	if isStackHeaderLine(line) {
		return true
	}
	if isTestDividerLine(line) {
		return true
	}
	if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "  ") {
		return true
	}
	if strings.Contains(line, ".go:") || strings.Contains(line, " +0x") {
		return true
	}
	if strings.HasPrefix(line, "created by ") || strings.HasPrefix(line, "panic"+"(") {
		return true
	}
	if strings.HasPrefix(line, "testing.") || strings.HasPrefix(line, "runtime.") {
		return true
	}
	if isPythonStackFrame(line) {
		return true
	}
	if isPythonExceptionOrError(line) {
		return true
	}
	if isAssertionDiffLine(line) {
		return true
	}

	return false
}

func isStackHeaderLine(line string) bool {
	if strings.Contains(line, "Traceback (most recent call last):") {
		return true
	}

	return strings.HasPrefix(line, "goroutine ") || strings.HasPrefix(line, "panic:")
}

func isTestDividerLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 10 {
		return false
	}

	return strings.HasPrefix(trimmed, "----------") || strings.HasPrefix(trimmed, "==========")
}

func isPythonStackFrame(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "File \"") && strings.Contains(trimmed, ", line ") {
		return true
	}

	return pyFrameRegex.CompileMust().MatchString(trimmed)
}

func isPythonExceptionOrError(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}
	if strings.HasPrefix(trimmed, "During handling of the above exception") ||
		strings.HasPrefix(trimmed, "The above exception was the direct cause") {
		return true
	}
	if pyExceptionRegex.CompileMust().MatchString(trimmed) {
		return true
	}
	colonIdx := strings.Index(trimmed, ":")
	if colonIdx > 0 {
		return isExceptionTypeName(strings.TrimSpace(trimmed[:colonIdx]))
	}

	return isExceptionTypeName(trimmed)
}

func isExceptionTypeName(name string) bool {
	if len(name) == 0 || strings.Contains(name, " ") {
		return false
	}
	suffixes := []string{"Error", "Exception", "Exit", "Interrupt", "Failure", "Warning"}
	for _, s := range suffixes {
		if strings.HasSuffix(name, s) {
			return true
		}
	}

	return false
}

func isAssertionDiffLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "--- expected") || strings.HasPrefix(trimmed, "+++ actual") {
		return true
	}
	if strings.HasPrefix(trimmed, "@@ ") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "+ ") || strings.HasPrefix(trimmed, "? ") {
		return true
	}

	return false
}

func isStackTerminator(line string) bool {
	if isStandardStackTerminator(line) {
		return true
	}
	if isRanTestsNoise(strings.TrimSpace(line)) || hasRanTestsNoisePattern(line) {
		return true
	}
	if isFailedSummaryNoise(line) {
		return true
	}
	if isGoTestDivider(line) {
		return true
	}

	return strings.Contains(line, "Process completed with exit code")
}

func isStandardStackTerminator(line string) bool {
	if line == "FAIL" || line == "PASS" || strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "PASS\t") {
		return true
	}
	if strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "?\t") {
		return true
	}

	return strings.HasPrefix(line, "##[") || strings.HasPrefix(line, "[command]")
}

func hasRanTestsNoisePattern(line string) bool {
	trimmed := strings.TrimSpace(line)

	return strings.HasPrefix(trimmed, "Ran ") && strings.Contains(trimmed, " tests in ")
}

func isFailedSummaryNoise(line string) bool {
	trimmed := strings.TrimSpace(line)

	return strings.HasPrefix(trimmed, "FAILED (failures=") || strings.HasPrefix(trimmed, "FAILED (errors=")
}

func isGoTestDivider(line string) bool {
	if strings.HasPrefix(line, "=== RUN") || strings.HasPrefix(line, "=== PAUSE") || strings.HasPrefix(line, "=== CONT") {
		return true
	}

	return strings.HasPrefix(line, "--- FAIL:") || strings.HasPrefix(line, "--- PASS:") || strings.HasPrefix(line, "--- SKIP:") || strings.HasPrefix(line, "--- BENCH:")
}

func filterCompactStackTrace(stack string) string {
	if len(stack) == 0 {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(stack), "\n")

	return strings.Join(capErrorLines(lines, 25), "\n")
}
