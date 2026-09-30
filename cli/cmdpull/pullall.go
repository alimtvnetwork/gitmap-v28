package cmdpull

import (
	"bytes"
	"strings"
)

// runPullAll is the explicit batch-pull entry point. It forwards to
// runPull with the --all flag injected, so the heavy lifting
// (resolution, parallelism, stop-on-fail, pending-task accounting)
// stays in one place. Useful for users and the right-click context
// menu's power-user "pull-all" action; equivalent to
// `gitmap pull --all <forwarded flags>`.
func runPullAll(args []string) error {
	return runPull(prependAll(args))
}

// RunPullAllJSON executes pull-all in-process with --json and returns the JSON string.
func RunPullAllJSON(args []string) (string, error) {
	var buf bytes.Buffer
	SetPullBatchJSONWriter(&buf)
	defer ResetPullBatchJSONWriter()

	jsonArgs := preparePullAllJSONArgs(args)
	err := runPull(jsonArgs)
	return buf.String(), err
}

func hasParallelFlag(args []string) bool {
	for _, a := range args {
		if a == "--parallel" || a == "-p" || strings.HasPrefix(a, "--parallel=") {
			return true
		}
	}
	return false
}

func stripJSONAndPullAllArgs(args []string) []string {
	var clean []string
	for _, a := range args {
		if strings.EqualFold(a, "--json") || strings.EqualFold(a, "-json") || a == "-j" {
			continue
		}
		if isPullAllFlagOrToken(a) {
			continue
		}
		clean = append(clean, a)
	}
	return clean
}

func preparePullAllJSONArgs(args []string) []string {
	clean := stripJSONAndPullAllArgs(args)
	base := []string{"--all", "--json"}
	if !hasParallelFlag(clean) {
		base = append(base, "--parallel", "2")
	}
	return append(base, clean...)
}

// prependAll injects --all at the front of args unless the caller
// already passed it (long or short form). Keeping it idempotent lets
// `gitmap pull-all --all` behave identically to `gitmap pull-all`.
func prependAll(args []string) []string {
	for _, a := range args {
		if isPullAllFlagOrToken(a) {
			return args
		}
	}

	return append([]string{"--all"}, args...)
}

func isPullAllFlagOrToken(a string) bool {
	return a == "--all" || a == "-all" || a == "-a" || a == "pat" || a == "pull-all-table"
}
