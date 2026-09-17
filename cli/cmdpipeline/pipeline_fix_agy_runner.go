package cmdpipeline

import (
	"strings"
)

// PipelineAgyFixRunner is an injectable runner for AGY pipeline fix integration.
var PipelineAgyFixRunner func(args []string) error

func isFixCompound(first string, rest []string) bool {
	return first == "fix" && hasAgyOrErrorsTarget(rest)
}

func isAgyCompound(first string, rest []string) bool {
	return first == "agy" && hasFixOrErrorsTarget(rest)
}

func isErrorLogsCompound(first string, rest []string) bool {
	return isErrorLogsSubcmd(first) && hasAgyOrFixTarget(rest)
}

// IsPipelineFixAgyArgs detects if args match pipeline fix agy / aef variations.
func IsPipelineFixAgyArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(args[0])
	rest := args[1:]
	if isDirectAgyFixToken(first) {
		return true
	}

	return isFixCompound(first, rest) || isAgyCompound(first, rest) || isErrorLogsCompound(first, rest)
}

func isDirectAgyFixToken(token string) bool {
	switch token {
	case "aef", "agy-errors-fix", "agy_errors_fix", "fix-errors-agy",
		"fix-agy", "fixagy", "pipeline-fix", "fix-pipeline":
		return true
	}

	return false
}

func hasAgyOrErrorsTarget(rest []string) bool {
	if len(rest) == 0 {
		return false
	}

	for _, arg := range rest {
		low := strings.ToLower(arg)
		if low == "agy" || low == "errors" || low == "error" || low == "aef" || low == "pipeline" {
			return true
		}
	}

	return false
}

func hasFixOrErrorsTarget(rest []string) bool {
	if len(rest) == 0 {
		return false
	}

	for _, arg := range rest {
		low := strings.ToLower(arg)
		if low == "fix" || low == "errors" || low == "error" || low == "fp" || low == "pipeline" {
			return true
		}
	}

	return false
}

func hasAgyOrFixTarget(rest []string) bool {
	if len(rest) == 0 {
		return false
	}

	for _, arg := range rest {
		low := strings.ToLower(arg)
		if low == "agy" || low == "fix" || low == "aef" {
			return true
		}
	}

	return false
}
