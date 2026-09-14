package cmd

import (
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

	return runGitPassthrough(append([]string{subCmd}, subArgs...))
}

func isPullAllGitSubCmd(subCmd string) bool {
	return subCmd == "pull-all" || subCmd == "pa"
}

func runGitPassthrough(args []string) error {
	if err := execGitInheritCP(args...); err != nil {
		return apperror.WrapSimple(err, "git passthrough failed")
	}

	return nil
}
