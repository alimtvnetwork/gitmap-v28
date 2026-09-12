package cmd

import (
	"os"
	"strconv"
)

// PipelineErrorFlags holds parsed flags and positional options for pipeline error-logs.
type PipelineErrorFlags struct {
	IsJSON            bool
	HasTimeline       bool
	HasFix            bool
	HasCheck          bool
	HasHelp           bool
	HasIndex          bool
	Index             int
	HasLastFailures   bool
	LastFailures      int
	HasLastFailedLogs bool
	FilePath          string
	TempFileName      string
}

// ParsePipelineErrorFlags parses command-line arguments for pipeline error-logs.
func ParsePipelineErrorFlags(args []string) PipelineErrorFlags {
	var flags PipelineErrorFlags
	parseCommonErrorFlags(args, &flags)
	parseIndexAndFailures(args, &flags)

	return flags
}

func parseCommonErrorFlags(args []string, flags *PipelineErrorFlags) {
	flags.HasHelp = hasArgFlag(args, "--help") || hasArgFlag(args, "-h")
	flags.IsJSON = hasArgFlag(args, "--json")
	flags.HasTimeline = hasTimelineArg(args)
	flags.HasFix = hasArgFlag(args, "--fix") || hasArgFlag(args, "-f")
	flags.HasCheck = hasArgFlag(args, "--check") || hasArgFlag(args, "-c")
	flags.FilePath = extractFlagVal(args, "--file")
	flags.TempFileName = extractFlagVal(args, "--tempfile")
}

func hasTimelineArg(args []string) bool {
	return hasArgFlag(args, "-t") || hasArgFlag(args, "--timeout") ||
		hasArgFlag(args, "--timeline") || hasArgFlag(args, "-w") || hasArgFlag(args, "--watch")
}

func parseIndexAndFailures(args []string, flags *PipelineErrorFlags) {
	flags.Index, flags.HasIndex = scanNegativeIndex(args)
	flags.LastFailures, flags.HasLastFailures = ParseLastFailuresFlag(args)
	flags.HasLastFailedLogs = hasArgFlag(args, "last-failed-logs") || hasArgFlag(os.Args, "last-failed-logs")
	if flags.HasLastFailedLogs && !flags.HasLastFailures {
		flags.LastFailures = 20
		flags.HasLastFailures = true
	}
}

// ParseNegativeIndex parses a string as a negative integer index (e.g. "-2").
func ParseNegativeIndex(arg string) (int, bool) {
	val, err := strconv.Atoi(arg)
	if err != nil || val >= 0 {
		return 0, false
	}

	return val, true
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
