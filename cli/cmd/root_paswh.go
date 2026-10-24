package cmd

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func runPASWH(args []string) error {
	if isPASWHHelp(args) {
		checkHelp("paswh", args)
		return nil
	}
	n, y, rest := parsePASWHArgs(args)
	pasArgs := append([]string{"--w", strconv.Itoa(n), "--hand", strconv.Itoa(y)}, rest...)

	return cmdssh.RunSSHPASFleet(pasArgs)
}

func parsePositiveIntArg(arg string, fallback int) int {
	parsed, err := strconv.Atoi(arg)
	if err == nil && parsed > 0 {
		return parsed
	}

	return fallback
}

func parsePASWHArgs(args []string) (int, int, []string) {
	n, y := 1, 1
	var rest []string
	posIdx := 0

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			rest = append(rest, arg)
			continue
		}
		switch posIdx {
		case 0:
			n = parsePositiveIntArg(arg, n)
			posIdx++
		case 1:
			y = parsePositiveIntArg(arg, y)
			posIdx++
		default:
			rest = append(rest, arg)
		}
	}

	return n, y, rest
}

func isPASWHHelp(args []string) bool {
	for i, a := range args {
		if a == "--help" || a == "help" {
			return true
		}
		if a == "-h" && !isNextArgNumeric(args, i) {
			return true
		}
	}

	return false
}

func isNextArgNumeric(args []string, i int) bool {
	if i+1 >= len(args) {
		return false
	}
	_, err := strconv.Atoi(args[i+1])

	return err == nil
}
