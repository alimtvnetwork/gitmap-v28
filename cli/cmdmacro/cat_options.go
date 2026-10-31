// Package cmdmacro — cat_options.go: argument parsing and configuration for gitmap cat.
package cmdmacro

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

// CatOptions holds parsed flags and parameters for file inspection.
type CatOptions struct {
	TargetFile     string
	StartLine      int
	EndLine        int
	HeadCount      int
	TailCount      int
	Limit          int
	Query          string
	HasLineNumbers bool
	HasLineRange   bool
	IsRaw          bool
}

// LineMatch holds an extracted line with its original 1-based line number.
type LineMatch struct {
	Number int
	Text   string
}

// DefaultCatOptions returns a zero-value initialized CatOptions with reasonable defaults.
func DefaultCatOptions() CatOptions {
	return CatOptions{
		StartLine: 1,
		EndLine:   -1,
		Limit:     50,
	}
}

// ParseCatOptions parses CLI arguments into a CatOptions struct.
func ParseCatOptions(args []string) (CatOptions, error) {
	opts := DefaultCatOptions()
	var rawTargets []string

	for idx := 0; idx < len(args); idx++ {
		arg := args[idx]
		if isCatFlag(arg) {
			idx = parseCatFlag(arg, args, idx, &opts)
			continue
		}

		rawTargets = append(rawTargets, arg)
	}

	opts.TargetFile = resolveRawCatTarget(rawTargets)

	return opts, nil
}

func isCatFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func resolveRawCatTarget(rawTargets []string) string {
	if len(rawTargets) == 0 {
		return ""
	}

	joined := strings.Join(rawTargets, " ")

	return macro.ExpandPathAndEnv(joined)
}

func parseCatFlag(arg string, args []string, idx int, opts *CatOptions) int {
	switch {
	case isLineRangeFlag(arg):
		return parseLineRangeFlag(args, idx, opts)
	case isHeadFlag(arg):
		return parseHeadFlag(args, idx, opts)
	case isTailFlag(arg):
		return parseTailFlag(args, idx, opts)
	case isQueryFlag(arg):
		return parseQueryFlag(args, idx, opts)
	case isLimitFlag(arg):
		return parseLimitFlag(args, idx, opts)
	case isLineNumberFlag(arg):
		return parseLineNumberFlag(args, idx, opts)
	case arg == "--raw":
		opts.IsRaw = true
		return idx
	default:
		return idx
	}
}

func isLineRangeFlag(arg string) bool {
	return arg == "-L" || arg == "--lines" || arg == "-lines"
}

func parseLineRangeFlag(args []string, idx int, opts *CatOptions) int {
	if idx+1 >= len(args) {
		return idx
	}

	val := args[idx+1]
	start, end := parseLineRange(val)
	opts.StartLine = start
	opts.EndLine = end
	opts.HasLineRange = true
	opts.HasLineNumbers = true

	return idx + 1
}

func parseLineRange(val string) (int, int) {
	sep := resolveRangeSeparator(val)
	if sep == "" {
		single, _ := strconv.Atoi(val)
		return single, single
	}

	parts := strings.SplitN(val, sep, 2)
	start, _ := strconv.Atoi(parts[0])
	if start <= 0 {
		start = 1
	}

	if parts[1] == "" {
		return start, -1
	}

	end, _ := strconv.Atoi(parts[1])

	return start, end
}

func resolveRangeSeparator(val string) string {
	if strings.Contains(val, ":") {
		return ":"
	}

	if strings.Contains(val, "-") && !strings.HasPrefix(val, "-") {
		return "-"
	}

	return ""
}

func isHeadFlag(arg string) bool {
	return arg == "--head" || arg == "-head"
}

func parseHeadFlag(args []string, idx int, opts *CatOptions) int {
	if idx+1 >= len(args) {
		return idx
	}

	cnt, err := strconv.Atoi(args[idx+1])
	if err == nil && cnt > 0 {
		opts.HeadCount = cnt
		return idx + 1
	}

	return idx
}

func isTailFlag(arg string) bool {
	return arg == "--tail" || arg == "-tail"
}

func parseTailFlag(args []string, idx int, opts *CatOptions) int {
	if idx+1 >= len(args) {
		return idx
	}

	cnt, err := strconv.Atoi(args[idx+1])
	if err == nil && cnt > 0 {
		opts.TailCount = cnt
		return idx + 1
	}

	return idx
}

func isQueryFlag(arg string) bool {
	return arg == "-q" || arg == "--query" || arg == "-s" || arg == "--filter"
}

func parseQueryFlag(args []string, idx int, opts *CatOptions) int {
	if idx+1 >= len(args) {
		return idx
	}

	opts.Query = args[idx+1]
	opts.HasLineNumbers = true

	return idx + 1
}

func isLimitFlag(arg string) bool {
	return arg == "-l" || arg == "--limit" || arg == "-limit"
}

func parseLimitFlag(args []string, idx int, opts *CatOptions) int {
	if idx+1 >= len(args) {
		return idx
	}

	cnt, err := strconv.Atoi(args[idx+1])
	if err == nil && cnt > 0 {
		opts.Limit = cnt
		return idx + 1
	}

	return idx
}

func isLineNumberFlag(arg string) bool {
	return arg == "--num" || arg == "-num" || arg == "--number" || arg == "--numbers" || arg == "-n"
}

func parseLineNumberFlag(args []string, idx int, opts *CatOptions) int {
	if argHasNextNumeric(args, idx) {
		cnt, _ := strconv.Atoi(args[idx+1])
		opts.HeadCount = cnt
		return idx + 1
	}

	opts.HasLineNumbers = true

	return idx
}

func argHasNextNumeric(args []string, idx int) bool {
	if idx+1 >= len(args) {
		return false
	}

	cnt, err := strconv.Atoi(args[idx+1])

	return err == nil && cnt > 0
}
