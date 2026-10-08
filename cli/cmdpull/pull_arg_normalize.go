package cmdpull

import (
	"strconv"
	"strings"
)

// NormalizePullArgs converts positional pull-all and table arguments into flags.
func NormalizePullArgs(args []string) []string {
	normalized := make([]string, 0, len(args)+1)
	for i := 0; i < len(args); i++ {
		if isPullAllTableSeq(args, i) {
			normalized = append(normalized, "--all", "--status")
			i++
			continue
		}
		if isHandShortFlagWithVal(args, i) {
			normalized = append(normalized, "--hand", args[i+1])
			i++
			continue
		}
		normalized = appendNormalizedPullToken(normalized, args[i])
	}

	return reorderPullFlags(normalized)
}

func isHandShortFlagWithVal(args []string, i int) bool {
	if args[i] != "-h" || i+1 >= len(args) {
		return false
	}
	_, err := strconv.Atoi(args[i+1])
	return err == nil
}

func reorderPullFlags(args []string) []string {
	flags := make([]string, 0, len(args))
	pos := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isPullFlagTakingValue(arg) {
			flags = appendValuedFlag(flags, args, &i)
			continue
		}
		flags, pos = categorizePullArg(arg, flags, pos)
	}

	return append(flags, pos...)
}

func categorizePullArg(arg string, flags, pos []string) ([]string, []string) {
	if strings.HasPrefix(arg, "-") {
		return append(flags, arg), pos
	}

	return flags, append(pos, arg)
}

func appendValuedFlag(flags []string, args []string, idx *int) []string {
	flags = append(flags, args[*idx])
	if *idx+1 < len(args) {
		*idx++
		flags = append(flags, args[*idx])
	}
	return flags
}

func isPullFlagTakingValue(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}

	return arg == "-g" || arg == "--group" || arg == "-p" || arg == "--parallel" ||
		arg == "-w" || arg == "--w" || arg == "--worker" || arg == "--workers" ||
		arg == "--hand" || arg == "--hands" || arg == "--h" || arg == "--concurrency"
}

func isPullAllTableSeq(args []string, i int) bool {
	isAllToken := args[i] == "all" || args[i] == "--all" || args[i] == "pa" || args[i] == "ta" || args[i] == "pull-all"

	return isAllToken && i+1 < len(args) && args[i+1] == "table"
}

func appendNormalizedPullToken(normalized []string, token string) []string {
	if isTableToken(token) {
		return append(normalized, "--all", "--status")
	}
	if isSSHFleetToken(token) {
		return append(normalized, "--all", "--ssh")
	}
	if isAllToken(token) {
		return append(normalized, "--all")
	}

	return append(normalized, token)
}

func isTableToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "pat" || lower == "pull-all-table"
}

func isAllToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "all" || lower == "pa" || lower == "ta" || lower == "pull-all"
}
