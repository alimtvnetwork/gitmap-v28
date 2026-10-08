package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
)

// runGitSubcommand handles transparent CLI routing for git subcommands.
func runGitSubcommand(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap git <command> [args...]", "E1000")
	}

	return dispatchGitSubcommand(args[0], args[1:])
}

func dispatchGitSubcommand(subCmd string, subArgs []string) error {
	if subCmd == "gitignore" || subCmd == "gitignore-agm" || subCmd == "agm" {
		return gitignoreagm.RunCLI(subArgs)
	}
	if subCmd == "ignore" || subCmd == "ig" {
		return cmdignore.RunIgnoreCLI(subArgs)
	}
	if subCmd == "pull" {
		return runPull(subArgs)
	}
	return dispatchGitPullSubcommands(subCmd, subArgs)
}

func dispatchGitPullSubcommands(subCmd string, subArgs []string) error {
	if subCmd == "pull-all-ssh" || subCmd == "pas" {
		return cmdpull.RunPullAll(append([]string{"--ssh"}, subArgs...))
	}
	if isPullAllGitSubCmd(subCmd) {
		return cmdpull.RunPullAll(subArgs)
	}
	if isPullAllTableGitSubCmd(subCmd) {
		return cmdpull.RunPullAll(append([]string{"--status"}, subArgs...))
	}

	return dispatchGitEfficientOrPassthrough(subCmd, subArgs)
}

func dispatchGitEfficientOrPassthrough(subCmd string, subArgs []string) error {
	if isPullEfficientTableGitSubCmd(subCmd) {
		return cmdpull.RunPullAllEfficient(subArgs, true, subCmd, isShortPullEfficientGitSubCmd(subCmd))
	}
	if isPullEfficientGitSubCmd(subCmd) {
		return cmdpull.RunPullAllEfficient(subArgs, false, subCmd, isShortPullEfficientGitSubCmd(subCmd))
	}

	return runGitPassthrough(append([]string{subCmd}, subArgs...))
}

func isPullAllGitSubCmd(subCmd string) bool {
	return subCmd == "pull-all" || subCmd == "pa" || subCmd == "ta"
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
	startTime := time.Now()
	err := cmdcommit.ExecGitInheritCP(args...)
	durationMs := time.Since(startTime).Milliseconds()
	exitCode := cmdcommit.ExtractGitExitCode(err)
	cmdcommit.RecordGitCommandHistory(args, exitCode, durationMs)
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "git passthrough failed")
	}

	return nil
}
