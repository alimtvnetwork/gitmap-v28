package cmdpipeline

import (
	"strconv"
	"strings"
)

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
	case "-f", "-c", "-v", "-V", "-t", "-w", "-h", "-n", "-y", "-j", "-ud":
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
	if flags.HasWorkers && val == strconv.Itoa(flags.Workers) {
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

// parseWorkersFlag parses --workers N (worker pool size override).
// Mirrors parseLimitFlag: accepts --workers=N or --workers N.
func parseWorkersFlag(args []string) (int, bool) {
	for i := 0; i < len(args); i++ {
		if val, isMatched := matchWorkersFlagToken(args, i); isMatched {
			return val, true
		}
	}

	return 0, false
}

func matchWorkersFlagToken(args []string, idx int) (int, bool) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--workers=") {
		return parsePositiveIntStr(strings.TrimPrefix(arg, "--workers="))
	}
	if arg != "--workers" || idx+1 >= len(args) {
		return 0, false
	}

	return parsePositiveIntStr(args[idx+1])
}
