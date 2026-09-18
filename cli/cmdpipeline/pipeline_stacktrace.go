package cmdpipeline

import (
	"bufio"
	"strings"
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
		job, _, text, _ := parseLogLine(raw)
		if isMatchingJobName(job, jobName) && len(text) > 0 {
			matched = append(matched, text)
		}
	}

	return strings.Join(matched, "\n")
}

func findStackTraceStart(raw string) int {
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

func extractBoundedStackLines(sub string) string {
	scanner := bufio.NewScanner(strings.NewReader(sub))
	var frames []string
	hasHeader := false
	for scanner.Scan() {
		if stop := processStackScanLine(scanner.Text(), &frames, &hasHeader); stop {
			break
		}
	}

	return strings.Join(frames, "\n")
}

func processStackScanLine(raw string, frames *[]string, hasHeader *bool) bool {
	line := sanitizeStackLine(raw)
	if *hasHeader && isStackStopLine(line, len(*frames)) {
		return true
	}
	if !*hasHeader && isStackStartLine(line) {
		*hasHeader = true
		*frames = append(*frames, line)
		return false
	}
	if !*hasHeader || len(line) == 0 {
		return false
	}
	*frames = append(*frames, line)

	return len(*frames) >= 25
}

func isStackStopLine(line string, frameCount int) bool {
	if len(line) == 0 && frameCount > 0 {
		return true
	}
	if isStackTerminator(line) {
		return true
	}

	return !isStackFrame(line)
}

func sanitizeStackLine(raw string) string {
	clean := ansiRegex.ReplaceAllString(raw, "")
	parts := strings.Split(clean, "\t")
	if len(parts) >= 3 {
		return cleanLogText(stripTimestamp(strings.Join(parts[2:], "\t")))
	}
	if len(parts) == 2 {
		return cleanLogText(parts[1])
	}

	return cleanLogText(clean)
}

func isStackStartLine(line string) bool {
	return strings.HasPrefix(line, "goroutine ") || strings.HasPrefix(line, "panic:") || strings.Contains(line, "panic:")
}

func isStackFrame(line string) bool {
	if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "  ") {
		return true
	}
	if strings.Contains(line, ".go:") || strings.Contains(line, " +0x") {
		return true
	}
	if strings.HasPrefix(line, "created by ") || strings.HasPrefix(line, "panic"+"(") {
		return true
	}

	return strings.HasPrefix(line, "testing.") || strings.HasPrefix(line, "runtime.")
}

func isStackTerminator(line string) bool {
	if strings.HasPrefix(line, "FAIL") || strings.HasPrefix(line, "PASS") {
		return true
	}
	if strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "?\t") {
		return true
	}
	if strings.HasPrefix(line, "##[") || strings.HasPrefix(line, "[command]") {
		return true
	}
	if strings.HasPrefix(line, "===") || strings.HasPrefix(line, "---") {
		return true
	}

	return strings.Contains(line, "Process completed with exit code")
}

func filterCompactStackTrace(stack string) string {
	if len(stack) == 0 {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(stack), "\n")

	return strings.Join(capErrorLines(lines, 25), "\n")
}
