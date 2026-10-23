package cmdpipeline

import (
	"os"
	"strconv"
	"strings"
)

// PipelineErrorFlags holds parsed flags and positional options for pipeline error-logs.
type PipelineErrorFlags struct {
	IsJSON               bool
	IsAll                bool
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
	HasLimit             bool
	Limit                int
	HasWorkers           bool
	Workers              int
	CommitTarget         string
	FilePath             string
	TempFileName         string
	HasForce             bool
	FormatProfile        string
	RepoTarget           string
	RawRepoTarget        string
	ResolvedPath         string
	HasAI                bool
	HasUntilDone         bool
}

// ParsePipelineErrorFlags parses command-line arguments for pipeline error-logs.
func ParsePipelineErrorFlags(args []string) PipelineErrorFlags {
	var flags PipelineErrorFlags
	parseCommonErrorFlags(args, &flags)
	parseIndexAndFailures(args, &flags)
	parseCommitTarget(args, &flags)
	parseRepoTarget(args, &flags)

	if flags.HasAI {
		SetGlobalAIFlag(true)
	}

	return flags
}

func parseCommonErrorFlags(args []string, flags *PipelineErrorFlags) {
	flags.HasHelp = hasArgFlag(args, "--help") || hasArgFlag(args, "-h") || hasArgFlag(args, "help")
	flags.IsJSON = hasArgFlag(args, "--json")
	flags.IsAll = hasArgFlag(args, "all") || hasArgFlag(args, "--all")
	flags.HasTimeline = hasTimelineArg(args)
	flags.HasCheck = hasArgFlag(args, "--check") || hasArgFlag(args, "-c")
	flags.IsDetailed = hasDetailedArg(args)
	flags.HasSuppressOutputLog = hasSuppressOutputArg(args)
	flags.Limit, flags.HasLimit = parseLimitFlag(args)
	flags.Workers, flags.HasWorkers = parseWorkersFlag(args)
	flags.FilePath = extractFlagVal(args, "--file")
	flags.TempFileName = extractFlagVal(args, "--tempfile")
	flags.HasForce = hasForceArg(args)
	flags.FormatProfile = extractFormatProfileFlag(args)
	flags.HasFix = hasArgFlag(args, "--fix") || (hasArgFlag(args, "-f") && flags.FormatProfile == "")
	flags.HasAI = hasArgFlag(args, "--ai") || hasArgFlag(os.Args, "--ai")
	flags.HasUntilDone = hasArgFlag(args, "-ud") || hasArgFlag(args, "--until-done") ||
		hasArgFlag(os.Args, "-ud") || hasArgFlag(os.Args, "--until-done")
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
	if hasArgFlag(args, "--no-output-log") || hasArgFlag(os.Args, "--no-output-log") {
		return true
	}

	return hasStandaloneShortNFlag(args) || hasStandaloneShortNFlag(os.Args)
}

func hasStandaloneShortNFlag(args []string) bool {
	for i, a := range args {
		if a == "-n" && !hasNextTokenPositiveInt(args, i) {
			return true
		}
	}

	return false
}

func hasNextTokenPositiveInt(args []string, idx int) bool {
	if idx+1 >= len(args) {
		return false
	}
	_, isPositive := parsePositiveIntStr(args[idx+1])

	return isPositive
}

func parseLimitFlag(args []string) (int, bool) {
	for i := 0; i < len(args); i++ {
		if val, isMatched := matchLimitFlagToken(args, i); isMatched {
			return val, true
		}
	}

	return 0, false
}

func matchLimitFlagToken(args []string, idx int) (int, bool) {
	arg := args[idx]
	if val, isInline := parseInlineLimit(arg); isInline {
		return val, true
	}
	if !isLimitFlagPrefix(arg) || idx+1 >= len(args) {
		return 0, false
	}

	return parsePositiveIntStr(args[idx+1])
}

func isLimitFlagPrefix(arg string) bool {
	switch {
	case strings.EqualFold(arg, "-l"),
		strings.EqualFold(arg, "-limit"),
		strings.EqualFold(arg, "--limit"),
		strings.EqualFold(arg, "-lines"),
		strings.EqualFold(arg, "--lines"),
		strings.EqualFold(arg, "-n"):
		return true
	default:
		return false
	}
}

func parseInlineLimit(arg string) (int, bool) {
	prefixes := []string{"--limit=", "-limit=", "--lines=", "-lines=", "-l=", "-n="}
	for _, p := range prefixes {
		if len(arg) >= len(p) && strings.EqualFold(arg[:len(p)], p) {
			return parsePositiveIntStr(arg[len(p):])
		}
	}

	return 0, false
}

func parsePositiveIntStr(raw string) (int, bool) {
	val, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || val <= 0 {
		return 0, false
	}

	return val, true
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

