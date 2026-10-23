package cmdpipeline

import (
	"strings"
)

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
