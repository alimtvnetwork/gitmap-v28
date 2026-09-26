package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// runGitSubcommand handles transparent CLI routing for git subcommands.
func runGitSubcommand(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap git <command> [args...]", "E1000")
	}

	return dispatchGitSubcommand(args[0], args[1:])
}

func dispatchGitSubcommand(subCmd string, subArgs []string) error {
	if subCmd == "pull" {
		return runPull(subArgs)
	}
	if isPullAllGitSubCmd(subCmd) {
		return runPullAll(subArgs)
	}
	if isPullAllTableGitSubCmd(subCmd) {
		return runPullAll(append([]string{"--status"}, subArgs...))
	}

	return dispatchGitEfficientOrPassthrough(subCmd, subArgs)
}

func dispatchGitEfficientOrPassthrough(subCmd string, subArgs []string) error {
	if isPullEfficientTableGitSubCmd(subCmd) {
		return runPullAllEfficient(subArgs, true, subCmd, isShortPullEfficientGitSubCmd(subCmd))
	}
	if isPullEfficientGitSubCmd(subCmd) {
		return runPullAllEfficient(subArgs, false, subCmd, isShortPullEfficientGitSubCmd(subCmd))
	}

	return runGitPassthrough(append([]string{subCmd}, subArgs...))
}

func isPullAllGitSubCmd(subCmd string) bool {
	return subCmd == "pull-all" || subCmd == "pa"
}

func isPullAllTableGitSubCmd(subCmd string) bool {
	lower := strings.ToLower(subCmd)

	return lower == "pull-all-table" || lower == "pat"
}

func isPullEfficientTableGitSubCmd(subCmd string) bool {
	lower := strings.ToLower(subCmd)

	return lower == "pull-all-efficient-table" || lower == "paet"
}

func isPullEfficientGitSubCmd(subCmd string) bool {
	lower := strings.ToLower(subCmd)

	return lower == "pull-all-efficient" || lower == "pae" || lower == "pull-ae"
}

func isShortPullEfficientGitSubCmd(subCmd string) bool {
	lower := strings.ToLower(subCmd)

	return lower == "pae" || lower == "pull-ae" || lower == "paet"
}

func runGitPassthrough(args []string) error {
	if err := execGitInheritCP(args...); err != nil {
		return apperror.WrapSimple(err, "git passthrough failed")
	}

	return nil
}
