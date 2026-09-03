package cmd

import (
	"bufio"
	"strings"
)

var failureMarkers = []string{
	"##[error]",
	"❌ FAIL",
	"--- FAIL:",
	"FAIL:",
	"FAILED",
	"fatal error:",
	"syntax error:",
	"exit status 1",
	"exit code 1",
	"panic:",
	"Stack Trace:",
}

// extractCleanErrorLines filters noisy logs to isolate only failure and error lines.
func extractCleanErrorLines(rawLogs string) string {
	if strings.TrimSpace(rawLogs) == "" {
		return ""
	}
	var matchedLines []string
	scanner := bufio.NewScanner(strings.NewReader(rawLogs))
	contextLines := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if hasFailureMarker(trimmed) {
			matchedLines = append(matchedLines, cleanMarkerLine(line))
			contextLines = 2
			continue
		}
		if contextLines > 0 {
			matchedLines = append(matchedLines, "    "+trimmed)
			contextLines--
		}
	}
	if len(matchedLines) == 0 {
		return extractTailLines(rawLogs, 20)
	}
	return strings.Join(matchedLines, "\n")
}

func hasFailureMarker(line string) bool {
	for _, m := range failureMarkers {
		if strings.Contains(line, m) {
			return true
		}
	}
	return false
}

func cleanMarkerLine(line string) string {
	line = strings.TrimPrefix(line, "##[error]")
	return strings.TrimSpace(line)
}

func extractTailLines(rawLogs string, n int) string {
	lines := strings.Split(strings.TrimSpace(rawLogs), "\n")
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}
