package cmdpipeline

import (
	"os"
	"strconv"
	"strings"
)

// PipelineErrorFlags holds parsed flags and positional options for pipeline error-logs.
type PipelineErrorFlags struct {
	IsJSON               bool
	HasTimeline          bool
	HasFix               bool
	HasCheck             bool
	HasHelp              bool
	HasIndex             bool
	Index                int
	HasLastFailures      bool
	LastFailures         int
	HasLastFailedLogs    bool
	IsDetailed           bool
	HasSuppressOutputLog bool
	CommitTarget         string
	FilePath             string
	TempFileName         string
	HasForce             bool
	FormatProfile        string
}

// ParsePipelineErrorFlags parses command-line arguments for pipeline error-logs.
func ParsePipelineErrorFlags(args []string) PipelineErrorFlags {
	var flags PipelineErrorFlags
	parseCommonErrorFlags(args, &flags)
	parseIndexAndFailures(args, &flags)
	parseCommitTarget(args, &flags)

	return flags
}

func parseCommonErrorFlags(args []string, flags *PipelineErrorFlags) {
	flags.HasHelp = hasArgFlag(args, "--help") || hasArgFlag(args, "-h") || hasArgFlag(args, "help")
	flags.IsJSON = hasArgFlag(args, "--json")
	flags.HasTimeline = hasTimelineArg(args)
	flags.HasCheck = hasArgFlag(args, "--check") || hasArgFlag(args, "-c")
	flags.IsDetailed = hasDetailedArg(args)
	flags.HasSuppressOutputLog = hasSuppressOutputArg(args)
	flags.FilePath = extractFlagVal(args, "--file")
	flags.TempFileName = extractFlagVal(args, "--tempfile")
	flags.HasForce = hasForceArg(args)
	flags.FormatProfile = extractFormatProfileFlag(args)
	flags.HasFix = hasArgFlag(args, "--fix") || (hasArgFlag(args, "-f") && flags.FormatProfile == "")
}

func extractFormatProfileFlag(args []string) string {
	val := extractFlagVal(args, "--format")
	if val != "" && !strings.HasPrefix(val, "-") {
		return val
	}
	val = extractFlagVal(args, "-f")
	if val != "" && !strings.HasPrefix(val, "-") {
		return val
	}

	return ""
}

func hasForceArg(args []string) bool {
	return hasArgFlag(args, "--force") || hasArgFlag(args, "--no-cache") ||
		hasArgFlag(os.Args, "--force") || hasArgFlag(os.Args, "--no-cache")
}

func hasSuppressOutputArg(args []string) bool {
	return hasArgFlag(args, "--no-output-log") || hasArgFlag(args, "-n") ||
		hasArgFlag(os.Args, "--no-output-log") || hasArgFlag(os.Args, "-n")
}

func hasDetailedArg(args []string) bool {
	return hasArgFlag(args, "--detailed") || hasArgFlag(args, "--verbose") ||
		hasArgFlag(args, "--v") || hasArgFlag(args, "-v") || hasArgFlag(args, "-V") ||
		hasArgFlag(os.Args, "--detailed") || hasArgFlag(os.Args, "--verbose") ||
		hasArgFlag(os.Args, "--v") || hasArgFlag(os.Args, "-v") || hasArgFlag(os.Args, "-V")
}

func hasTimelineArg(args []string) bool {
	return hasArgFlag(args, "-t") || hasArgFlag(args, "--timeout") ||
		hasArgFlag(args, "--timeline") || hasArgFlag(args, "-w") || hasArgFlag(args, "--watch")
}

func parseIndexAndFailures(args []string, flags *PipelineErrorFlags) {
	flags.Index, flags.HasIndex = scanNegativeIndex(args)
	flags.LastFailures, flags.HasLastFailures = ParseLastFailuresFlag(args)
	flags.HasLastFailedLogs = hasArgFlag(args, "last-failed-logs") || hasArgFlag(os.Args, "last-failed-logs")
	if flags.HasLastFailedLogs && flags.LastFailures == 0 {
		flags.LastFailures = 20
		flags.HasLastFailures = true
	}
}

// ParseNegativeIndex parses a string as a negative integer index (e.g. "-2", "-1n", "HEAD~1").
func ParseNegativeIndex(arg string) (int, bool) {
	trimmed := strings.TrimSpace(arg)
	if offset, isHead := parseHeadRevisionOffset(trimmed); isHead {
		return offset, true
	}

	clean := strings.TrimSuffix(strings.TrimSuffix(trimmed, "n"), "N")
	val, err := strconv.Atoi(clean)
	if err != nil || val >= 0 {
		return 0, false
	}

	return val, true
}

func parseHeadRevisionOffset(arg string) (int, bool) {
	lower := strings.ToLower(arg)
	isHeadTilde := strings.HasPrefix(lower, "head~")
	isTildeOnly := strings.HasPrefix(lower, "~")
	if !isHeadTilde && !isTildeOnly {
		return 0, false
	}

	raw := strings.TrimPrefix(strings.TrimPrefix(lower, "head~"), "~")
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return 0, false
	}

	return -val, true
}

// IsNegativeIndexToken reports whether arg is a negative integer token.
func IsNegativeIndexToken(arg string) bool {
	_, hasValidIndex := ParseNegativeIndex(arg)

	return hasValidIndex
}

func scanNegativeIndex(args []string) (int, bool) {
	for _, arg := range args {
		if val, hasIndex := ParseNegativeIndex(arg); hasIndex {
			return val, true
		}
	}

	return 0, false
}

// NormalizeNegativeIndex converts a negative 1-based index (e.g. -1, -2) to a 0-based slice index.
func NormalizeNegativeIndex(offset int) int {
	if offset < 0 {
		return -offset - 1
	}

	return offset
}

// NormalizeCommitOffset converts a commit offset (0, -1, -2) to a 0-based commit history index.
func NormalizeCommitOffset(offset int) int {
	if offset < 0 {
		return -offset
	}

	return offset
}

// ParseLastFailuresFlag parses --last-failures <N> or --last-failures=<N>.
func ParseLastFailuresFlag(args []string) (int, bool) {
	raw := extractFlagVal(args, "--last-failures")
	if len(raw) == 0 {
		return 0, false
	}

	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return 0, false
	}

	return val, true
}

func parseCommitTarget(args []string, flags *PipelineErrorFlags) {
	explicit := extractFlagVal(args, "--commit")
	if len(explicit) > 0 {
		flags.CommitTarget = explicit

		return
	}

	flags.CommitTarget = scanCommitTargetOrOffset(args)
}

func scanCommitTargetOrOffset(args []string) string {
	for _, arg := range args {
		target := evaluateCommitCandidate(arg)
		if len(target) > 0 {
			return target
		}
	}

	return ""
}

func evaluateCommitCandidate(arg string) string {
	trimmed := strings.TrimSpace(arg)
	if isSkipTokenForCommit(trimmed) {
		return ""
	}
	if _, isOffset := ParseNegativeIndex(trimmed); isOffset {
		return trimmed
	}
	if isCommitHexSha(trimmed) {
		return trimmed
	}
	lower := strings.ToLower(trimmed)
	if lower == "latest" || lower == "head" {
		return trimmed
	}

	return ""
}

func isSkipTokenForCommit(token string) bool {
	if len(token) == 0 || strings.HasPrefix(token, "--") {
		return true
	}

	switch token {
	case "-f", "-c", "-v", "-V", "-t", "-w", "-h", "-n", "-y", "-j":
		return true
	case "clear", "last-failed-logs", "errors", "error-logs", "pe":
		return true
	default:
		return false
	}
}

func isCommitHexSha(s string) bool {
	if strings.HasPrefix(s, "-") || len(s) < 4 || len(s) > 40 {
		return false
	}

	for _, r := range strings.ToLower(s) {
		isDigit := r >= '0' && r <= '9'
		isHexLetter := r >= 'a' && r <= 'f'
		if !isDigit && !isHexLetter {
			return false
		}
	}

	return true
}
