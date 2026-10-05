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
	CommitTarget         string
	FilePath             string
	TempFileName         string
	HasForce             bool
	FormatProfile        string
	RepoTarget           string
	RawRepoTarget        string
	ResolvedPath         string
}

// ParsePipelineErrorFlags parses command-line arguments for pipeline error-logs.
func ParsePipelineErrorFlags(args []string) PipelineErrorFlags {
	var flags PipelineErrorFlags
	parseCommonErrorFlags(args, &flags)
	parseIndexAndFailures(args, &flags)
	parseCommitTarget(args, &flags)
	parseRepoTarget(args, &flags)

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
	case "clear", "last-failed-logs", "errors", "error-logs", "pe", "history-ai", "hai":
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

func parseRepoTarget(args []string, flags *PipelineErrorFlags) {
	if applyExplicitRepoTarget(args, flags) {
		return
	}

	findPositionalRepoTarget(args, flags)
}

func applyExplicitRepoTarget(args []string, flags *PipelineErrorFlags) bool {
	explicit := extractFlagVal(args, "--repo")
	if explicit == "" {
		explicit = extractFlagVal(args, "-r")
	}
	if len(explicit) == 0 {
		return false
	}
	flags.RawRepoTarget = explicit
	flags.RepoTarget, flags.ResolvedPath = ResolvePipelineTargetAndPath(explicit)

	return true
}

func findPositionalRepoTarget(args []string, flags *PipelineErrorFlags) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if isValueFlag(a) {
			i += advanceValueFlagOffset(a, flags.HasLimit)
			continue
		}
		if applyCandidateRepoTarget(a, flags) {
			return
		}
	}
}

func advanceValueFlagOffset(flagName string, hasLimit bool) int {
	if flagName != "-n" || hasLimit {
		return 1
	}

	return 0
}

func applyCandidateRepoTarget(arg string, flags *PipelineErrorFlags) bool {
	trimmed := strings.TrimSpace(arg)

	if strings.EqualFold(trimmed, "all") {
		flags.IsAll = true

		return true
	}

	if isSkipTokenForRepoTarget(trimmed) || trimmed == flags.CommitTarget || isConsumedFlagValue(trimmed, flags) {
		return false
	}
	flags.RawRepoTarget = trimmed
	flags.RepoTarget, flags.ResolvedPath = ResolvePipelineTargetAndPath(trimmed)

	return true
}

func isValueFlag(flag string) bool {
	switch strings.ToLower(flag) {
	case "--file", "--tempfile", "--format", "--repo", "-r", "--last-failures",
		"-l", "-limit", "--limit", "--lines", "-lines", "-n":
		return true
	default:
		return false
	}
}

func isConsumedFlagValue(val string, flags *PipelineErrorFlags) bool {
	if flags.HasLimit && val == strconv.Itoa(flags.Limit) {
		return true
	}

	return val == flags.FilePath || val == flags.TempFileName || val == flags.FormatProfile
}

func isSkipTokenForRepoTarget(token string) bool {
	if len(token) == 0 || strings.HasPrefix(token, "-") {
		return true
	}
	switch strings.ToLower(token) {
	case "clear", "last-failed-logs", "errors", "error-logs", "pe", "ee", "help",
		"format", "add-format", "rm-format", "remove-format", "add-all", "list-formats", "preview-format",
		"history-ai", "hai", "all":
		return true
	default:
		return false
	}
}
