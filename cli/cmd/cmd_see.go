package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

// HistoryRunnerFn is a hook to delegate history execution from cmd package.
var HistoryRunnerFn func(args []string) *apperror.AppError

// RunSeeCLI routes gitmap see subcommands.
func RunSeeCLI(args []string) *apperror.AppError {
	isZero := len(args) == 0
	if isZero {
		return PrintSeeHelp()
	}

	joined := strings.ToLower(strings.Join(args, " "))

	return routeSeeCommand(args, joined)
}

func routeSeeCommand(args []string, joined string) *apperror.AppError {
	if isCommitPendingCmd(joined) {
		return runSeeCommitPending()
	}

	if hasIgnorePrefix(joined) {
		return routeIgnoreIssues(args)
	}

	if hasErrorsSSHPrefix(joined) {
		return runSeeErrorsSSH(args)
	}

	if hasHistorySSHPrefix(joined) {
		return runSeeHistorySSH(args)
	}

	return routeRemaining(args, joined)
}

func isCommitPendingCmd(joined string) bool {
	return strings.HasPrefix(joined, "commit pending") ||
		strings.HasPrefix(joined, "commit-pending")
}

func hasIgnorePrefix(joined string) bool {
	return strings.HasPrefix(joined, "git-ignore") ||
		strings.HasPrefix(joined, "ignore") ||
		strings.HasPrefix(joined, "ig")
}

func hasErrorsSSHPrefix(joined string) bool {
	return strings.HasPrefix(joined, "errors ssh") ||
		strings.HasPrefix(joined, "ses") ||
		joined == "ses"
}

func hasHistorySSHPrefix(joined string) bool {
	return strings.HasPrefix(joined, "history ssh") ||
		strings.HasPrefix(joined, "nodes history") ||
		strings.HasPrefix(joined, "nodes histories")
}

func routeRemaining(args []string, joined string) *apperror.AppError {
	if strings.HasPrefix(joined, "errors") {
		err := cmderrors.RunErrorsCLI(args[1:])

		return apperror.WrapSimple(err, "errors")
	}

	if strings.HasPrefix(joined, "history") {
		return runSeeHistory(args[1:])
	}

	if strings.HasPrefix(joined, "repo-manage") || strings.HasPrefix(joined, "ui") {
		return RunRepoManageUI()
	}

	return apperror.NewSimple("unknown see command: "+joined, "E1020")
}

func routeIgnoreIssues(args []string) *apperror.AppError {
	var scanArgs []string
	isLonger := len(args) > 1
	if isLonger {
		scanArgs = filterIgnoreFlags(args[1:])
	}

	err := cmdignore.RunIgnoreCLI(append([]string{"scan"}, scanArgs...))

	return apperror.WrapSimple(err, "ignore")
}

func filterIgnoreFlags(subArgs []string) []string {
	var filtered []string
	for _, a := range subArgs {
		isIssueWord := strutil.EqualFoldAny(a, "issues", "issue")
		if !isIssueWord {
			filtered = append(filtered, a)
		}
	}

	return filtered
}
