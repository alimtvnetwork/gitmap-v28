package cmdgit

import (
	"slices"
	"strings"
)

// RunLog handles native git log execution with clean, modern output.
func RunLog(args []string) error {
	targetDir, remaining := ParseGitDirAndArgs(args)

	gitArgs := buildLogArgs(remaining)

	return ExecGitInspection(targetDir, gitArgs...)
}

func buildLogArgs(args []string) []string {
	var finalArgs []string
	finalArgs = append(finalArgs, "log")

	hasLimit := hasLimitFlag(args)
	if !hasLimit {
		finalArgs = append(finalArgs, "-n", "10")
	}

	hasCustomFormat := hasFormatOrDetailFlag(args)
	if !hasCustomFormat {
		finalArgs = append(finalArgs,
			"--date=short",
			"--pretty=format:%C(yellow)%h%C(reset) %C(green)%cd%C(reset) %C(cyan)%an%C(reset) %s%C(auto)%d%C(reset)",
		)
	}

	finalArgs = append(finalArgs, args...)

	return finalArgs
}

func hasLimitFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-n" || strings.HasPrefix(arg, "-n") || strings.HasPrefix(arg, "--max-count") {
			return true
		}

		if len(arg) > 1 && arg[0] == '-' && isAllDigits(arg[1:]) {
			return true
		}
	}

	return false
}

var logFormatFlagTokens = []string{"--stat", "--shortstat", "-p", "--patch", "--oneline"}

func hasFormatOrDetailFlag(args []string) bool {
	for _, arg := range args {
		if slices.Contains(logFormatFlagTokens, arg) {
			return true
		}

		if strings.HasPrefix(arg, "--pretty") || strings.HasPrefix(arg, "--format") {
			return true
		}
	}

	return false
}

func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
