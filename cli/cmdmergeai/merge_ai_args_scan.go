// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"strings"
)

type parsedFlagState struct {
	dest        string
	sources     []string
	configPath  string
	positionals []string
	dryRun      bool
	isHelp      bool
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		lower := strings.ToLower(strings.TrimSpace(a))
		if lower == "-h" || lower == "--help" || lower == "help" {
			return true
		}
	}
	return false
}

func isDestFlag(lower string) bool {
	return lower == "-d" || lower == "--destination" || lower == "--dest" || lower == "--target"
}

func isDestAssignFlag(lower string) bool {
	return strings.HasPrefix(lower, "--destination=") || strings.HasPrefix(lower, "-d=") ||
		strings.HasPrefix(lower, "--dest=") || strings.HasPrefix(lower, "--target=")
}

func isSourceFlag(lower string) bool {
	return lower == "-s" || lower == "--source" || lower == "-u" || lower == "--url"
}

func isSourceAssignFlag(lower string) bool {
	return strings.HasPrefix(lower, "--source=") || strings.HasPrefix(lower, "-s=") ||
		strings.HasPrefix(lower, "--url=") || strings.HasPrefix(lower, "-u=")
}

func isConfigFlag(lower string) bool {
	return lower == "-c" || lower == "--config"
}

func isConfigAssignFlag(lower string) bool {
	return strings.HasPrefix(lower, "--config=") || strings.HasPrefix(lower, "-c=")
}

func extractAssignValue(arg string) string {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func peekNextArg(args []string, i int) (string, bool) {
	if i+1 < len(args) {
		return strings.TrimSpace(args[i+1]), true
	}
	return "", false
}

func tryScanArg(args []string, i int, s *parsedFlagState) (int, bool) {
	a := strings.TrimSpace(args[i])
	lower := strings.ToLower(a)

	if lower == "-h" || lower == "--help" || lower == "help" {
		s.isHelp = true
		return i + 1, true
	}

	if lower == "--dry-run" {
		s.dryRun = true
		return i + 1, true
	}

	if isDestAssignFlag(lower) {
		s.dest = extractAssignValue(a)
		return i + 1, true
	}

	if isSourceAssignFlag(lower) {
		s.sources = append(s.sources, extractAssignValue(a))
		return i + 1, true
	}

	if isConfigAssignFlag(lower) {
		s.configPath = extractAssignValue(a)
		return i + 1, true
	}

	nextVal, hasNext := peekNextArg(args, i)

	if isDestFlag(lower) && hasNext {
		s.dest = nextVal
		return i + 2, true
	}

	if isSourceFlag(lower) && hasNext {
		s.sources = append(s.sources, nextVal)
		return i + 2, true
	}

	if isConfigFlag(lower) && hasNext {
		s.configPath = nextVal
		return i + 2, true
	}

	return i, false
}

func scanArgsIntoState(args []string) parsedFlagState {
	var s parsedFlagState
	i := 0
	for i < len(args) {
		nextIdx, handled := tryScanArg(args, i, &s)
		if handled {
			i = nextIdx
			continue
		}
		s.positionals = append(s.positionals, strings.TrimSpace(args[i]))
		i++
	}
	return s
}
